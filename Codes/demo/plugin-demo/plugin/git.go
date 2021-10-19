package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitDiff struct {
	Action          string `yaml:"action" json:"action" bson:"action" validate:""`
	FileName        string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Stats           string `yaml:"stats" json:"stats" bson:"stats" validate:""`
	ModifyLineCount int    `yaml:"modifyLineCount" json:"modifyLineCount" bson:"modifyLineCount" validate:""`
}

type Commit struct {
	ProjectName          string    `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	RunName              string    `yaml:"runName" json:"runName" bson:"runName" validate:""`
	BranchName           string    `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Commit               string    `yaml:"commit" json:"commit" bson:"commit" validate:""`
	FullMessage          string    `yaml:"fullMessage" json:"fullMessage" bson:"fullMessage" validate:""`
	CommitTime           time.Time `yaml:"commitTime" json:"commitTime" bson:"commitTime" validate:""`
	CommitterName        string    `yaml:"committerName" json:"committerName" bson:"committerName" validate:""`
	CommitterEmail       string    `yaml:"committerEmail" json:"committerEmail" bson:"committerEmail" validate:""`
	Message              string    `yaml:"message" json:"message" bson:"message" validate:""`
	Diffs                []GitDiff `yaml:"diffs" json:"diffs" bson:"diffs" validate:""`
	ModifyLineCountTotal int       `yaml:"modifyLineCountTotal" json:"modifyLineCountTotal" bson:"modifyLineCountTotal" validate:""`
	CreateTime           time.Time `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
}

type Git struct {
}

func gitGetRefsAndBranchHead(repo *git.Repository, branch string) ([]string, plumbing.Hash, bool, error) {
	var err error
	var strRefs []string
	var branchHead plumbing.Hash
	needCheckout := true

	refs, err := repo.References()
	if err != nil {
		return strRefs, branchHead, needCheckout, err
	}
	defer refs.Close()

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}
		if ref.Name().String() == fmt.Sprintf("refs/remotes/origin/%s", branch) {
			branchHead = ref.Hash()
		}
		if ref.Name().String() == fmt.Sprintf("refs/heads/%s", branch) {
			needCheckout = false
		}
		strRefs = append(strRefs, ref.String())
		return nil
	})
	if err != nil {
		return strRefs, branchHead, needCheckout, err
	}
	return strRefs, branchHead, needCheckout, err
}

func printStat(fileStats []object.FileStat) string {
	padLength := float64(len(" "))
	newlineLength := float64(len("\n"))
	separatorLength := float64(len("|"))
	// Soft line length limit. The text length calculation below excludes
	// length of the change number. Adding that would take it closer to 80,
	// but probably not more than 80, until it's a huge number.
	lineLength := 72.0

	// Get the longest filename and longest total change.
	var longestLength float64
	var longestTotalChange float64
	for _, fs := range fileStats {
		if int(longestLength) < len(fs.Name) {
			longestLength = float64(len(fs.Name))
		}
		totalChange := fs.Addition + fs.Deletion
		if int(longestTotalChange) < totalChange {
			longestTotalChange = float64(totalChange)
		}
	}

	// Parts of the output:
	// <pad><filename><pad>|<pad><changeNumber><pad><+++/---><newline>
	// example: " main.go | 10 +++++++--- "

	// <pad><filename><pad>
	leftTextLength := padLength + longestLength + padLength

	// <pad><number><pad><+++++/-----><newline>
	// Excluding number length here.
	rightTextLength := padLength + padLength + newlineLength

	totalTextArea := leftTextLength + separatorLength + rightTextLength
	heightOfHistogram := lineLength - totalTextArea

	// Scale the histogram.
	var scaleFactor float64
	if longestTotalChange > heightOfHistogram {
		// Scale down to heightOfHistogram.
		scaleFactor = longestTotalChange / heightOfHistogram
	} else {
		scaleFactor = 1.0
	}

	finalOutput := ""
	for _, fs := range fileStats {
		addn := float64(fs.Addition)
		deln := float64(fs.Deletion)
		addCount := int(math.Abs(math.Floor(addn / scaleFactor)))
		delCount := int(math.Abs(math.Floor(deln / scaleFactor)))
		adds := strings.Repeat("+", addCount)
		dels := strings.Repeat("-", delCount)
		finalOutput += fmt.Sprintf("%d %s%s", fs.Addition+fs.Deletion, adds, dels)
	}

	return finalOutput
}

func gitGetDiff(repo *git.Repository, previousCommit string) ([]Commit, error) {
	var err error
	var gcs []Commit
	commits, err := repo.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		return gcs, err
	}
	defer commits.Close()

	var commitArray []*object.Commit
	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}
		commitArray = append([]*object.Commit{commit}, commitArray...)
		if commit.Hash.String() == previousCommit {
			break
		}
	}

	var prevCommit *object.Commit
	var prevTree *object.Tree

	for _, commit := range commitArray {
		currentTree, err := commit.Tree()
		if err != nil {
			return gcs, err
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		changes, err := currentTree.Diff(prevTree)
		if err != nil {
			return gcs, err
		}

		var gds []GitDiff
		modifyLineCountTotal := 0
		for _, c := range changes {
			action, err := c.Action()
			var strAction, strFromFile, strToFile, strPatch string
			if err == nil {
				strAction = action.String()
			}
			ffs, tfs, err := c.Files()
			if err == nil {
				if ffs != nil {
					strFromFile = ffs.Name
				}
				if tfs != nil {
					strToFile = tfs.Name
				}
			}
			if strFromFile != strToFile && strFromFile != "" && strToFile != "" {
				strAction = "Rename"
			}
			patch, _ := c.Patch()
			modifyLineCount := 0
			var fileName string
			if err == nil {
				for _, stat := range patch.Stats() {
					modifyLineCount = modifyLineCount + stat.Addition + stat.Deletion
					fileName = stat.Name
				}
				strPatch = printStat(patch.Stats())
			}
			gd := GitDiff{
				Action:          strAction,
				FileName:        fileName,
				Stats:           strPatch,
				ModifyLineCount: modifyLineCount,
			}
			modifyLineCountTotal = modifyLineCountTotal + modifyLineCount
			gds = append(gds, gd)
		}
		gc := Commit{
			Commit:               commit.Hash.String(),
			FullMessage:          commit.String(),
			CommitTime:           commit.Committer.When,
			CommitterName:        commit.Committer.Name,
			CommitterEmail:       commit.Committer.Email,
			Message:              commit.Message,
			Diffs:                gds,
			ModifyLineCountTotal: modifyLineCountTotal,
		}
		gcs = append(gcs, gc)

		prevCommit = commit
		prevTree = currentTree
	}
	return gcs, err
}

func gitGetLatestTag(repo *git.Repository, branchHead plumbing.Hash) string {
	var tagName string
	var err error
	// get all tags
	tags, err := repo.Tags()
	if err != nil {
		panic(err)
	}
	defer tags.Close()
	tagsMap := map[plumbing.Hash]*plumbing.Reference{}
	_ = tags.ForEach(func(tag *plumbing.Reference) error {
		tagsMap[tag.Hash()] = tag
		return nil
	})

	commits, err := repo.Log(&git.LogOptions{
		From: branchHead,
	})
	if err != nil {
		return tagName
	}
	defer commits.Close()

	var tag *plumbing.Reference
	var count int
	_ = commits.ForEach(func(c *object.Commit) error {
		if t, ok := tagsMap[c.Hash]; ok {
			tag = t
		}
		if tag != nil {
			return storer.ErrStop
		}
		count++
		return nil
	})
	if tag != nil {
		if count == 0 {
			tagName = tag.Name().Short()
		} else {
			tagName = fmt.Sprintf("%v-%v-g%v", tag.Name().Short(), count, branchHead.String()[0:7])
		}
	}
	return tagName
}

func gitCheckout(repo *git.Repository, branch string, branchHead plumbing.Hash, needCheckout bool) error {
	var err error
	if !branchHead.IsZero() && needCheckout {
		w, err := repo.Worktree()
		if err != nil {
			return err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Hash:   branchHead,
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
		})
		if err != nil {
			return err
		}
	}
	return err
}

// GitPull: clone repository from git and check out to branch
func (g Git) GitPull(dir, url, branch, username, password string, timeoutSeconds int, previousCommit string) ([]Commit, string, string, error) {
	errInfo := fmt.Sprintf("git clone %s error", url)
	var err error
	var latestCommit, tagName string
	gcs := []Commit{}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		_ = os.RemoveAll(dir)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = fmt.Errorf("%s: mkdir %s error: %s", errInfo, dir, err.Error())
		fmt.Println("[ERROR]", err.Error())
		return gcs, latestCommit, tagName, err
	}

	customClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Duration(timeoutSeconds) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	client.InstallProtocol("https", githttp.NewClient(customClient))

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		// git clone
		var r *git.Repository
		r, err = git.PlainClone(dir, false, &git.CloneOptions{
			Auth: &githttp.BasicAuth{
				Username: username,
				Password: password,
			},
			URL:      url,
			Progress: writer,
		})
		if err != nil {
			err = fmt.Errorf("%s: clone repository error: %s", errInfo, err.Error())
			fmt.Println("[ERROR]", err.Error())
			return
		}

		var strRefs []string
		var branchHead plumbing.Hash
		var needCheckout bool
		strRefs, branchHead, needCheckout, err = gitGetRefsAndBranchHead(r, branch)
		if err != nil {
			err = fmt.Errorf("%s: get refs error: %s", errInfo, err.Error())
			fmt.Println("[ERROR]", err.Error())
			return
		}
		fmt.Println("[INFO]", "git show-ref")
		for _, strRef := range strRefs {
			fmt.Println("[INFO]", strRef)
		}

		latestCommit = branchHead.String()
		tagName = gitGetLatestTag(r, branchHead)

		gcs, err = gitGetDiff(r, previousCommit)
		if err != nil {
			err = fmt.Errorf("%s: get commits diff error: %s", errInfo, err.Error())
			fmt.Println("[ERROR]", err.Error())
			return
		}
		fmt.Println("[INFO]", "git commit diff")
		fmt.Println("[INFO]", "####################################")
		for _, gc := range gcs {
			scanner := bufio.NewScanner(bytes.NewBuffer([]byte(gc.FullMessage)))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				fmt.Println("[INFO]", scanner.Text())
			}
			for _, diff := range gc.Diffs {
				var msg string
				msg = fmt.Sprintf("%s %s | %s", diff.Action, diff.FileName, diff.Stats)
				fmt.Println("[INFO]", msg)
			}
			fmt.Println("[INFO]", "####################################")
		}

		// checkout branch
		err = gitCheckout(r, branch, branchHead, needCheckout)
		if err != nil {
			err = fmt.Errorf("%s: git checkout error: %s", errInfo, err.Error())
			fmt.Println("[ERROR]", err.Error())
			return
		}

		fmt.Println("[INFO]", fmt.Sprintf("tagName: %s", tagName))
		fmt.Println("[INFO]", fmt.Sprintf("branch %s HEAD hash: %v", branch, branchHead.String()))
	}()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println("[INFO]", scanner.Text())
	}
	return gcs, latestCommit, tagName, err
}

func init() {
	fmt.Println("git plugin init")
}

var GitPlugin Git

// go build -o plugin.so -buildmode=plugin git.go
