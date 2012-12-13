package svn

import "github.com/kless/goconfig/config"
import "os/exec"
import "log"
import "strings"

type Svn struct {
	LastErr string
	Config *config.Config
	blocked bool
}

func (s *Svn) getBinPath() (string) {
	val,_ := s.Config.String("SVN", "binpath")
	return val
}

func (s *Svn) getCheckoutPath() (string) {
	val,_ := s.Config.String("SVN", "checkout")
	return val
}

func (s *Svn) getRepositoryPath() (string) {
	val,_ := s.Config.String("SVN", "repository")
	return val
}

func (s *Svn) shouldAlwaysRevert() (bool) {
	val, _ := s.Config.Bool("CMD", "always_revert_first")
	return val
}

func (s *Svn) SvnLocalInfo() (*exec.Cmd){
	cmd := exec.Command(s.getBinPath(), "info", s.getCheckoutPath())
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) SvnRemoteInfo() (*exec.Cmd) {
	cmd := exec.Command(s.getBinPath(), "info", s.getCheckoutPath())
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) UpdateWorkingCopy() (*exec.Cmd) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	cmd := exec.Command(s.getBinPath(), "up")
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) SwitchTrunk() (*exec.Cmd) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/trunk")
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) SwitchBranch(branchName string) (*exec.Cmd) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	log.Printf("%s %s %s", s.getBinPath(), "sw", s.getRepositoryPath() + "/branches/" + branchName)
	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/branches/" + branchName)
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) SwitchTag(tagName string) (*exec.Cmd) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/tags/" + tagName)
	cmd.Dir = s.getCheckoutPath()
	return cmd
}

func (s *Svn) getRemoteDirList(pathInRepo string) ([]string) {
	cmd := exec.Command(s.getBinPath(), "ls", s.getRepositoryPath() + "/" + pathInRepo)
	cmd.Dir = s.getCheckoutPath()
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Panic(err)
	}
	outstring := strings.TrimSpace(string(out))
	dirList := strings.Split(outstring, "\n")

	return dirList
}

func (s *Svn) GetBranchList() ([]string) {
	return s.getRemoteDirList("branches")
}

func (s *Svn) GetTagList() ([]string) {
	return s.getRemoteDirList("tags")
}

func (s *Svn) revertWorkingCopy() {
	cmd := exec.Command(s.getBinPath(), "revert", "-R", ".")
	cmd.Dir = s.getCheckoutPath()
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(out)
		log.Panic(err)
	}
}
