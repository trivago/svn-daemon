package vcs

import "github.com/kless/goconfig/config"
import "os/exec"
import "log"
import "strings"

type Svn struct {
	LastErr string
	Config *config.Config
	blocked bool
}

func (s *Svn) getHookScriptResult(hookname string) (string) {
	val, _ := s.Config.String("CMD", hookname)
	if val + "FOOBAR" == "FOOBAR" {
		return ""
	}

	workdir, _ := s.Config.String("SVN", "checkout")
	cmd := exec.Command(val)
	cmd.Dir = workdir
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(err.Error())
	}

	return string(out)
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

func (s *Svn) UpdateWorkingCopy() (string) {

	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	preOut := s.getHookScriptResult("pre_up")
	cmd := exec.Command(s.getBinPath(), "up")
	cmd.Dir = s.getCheckoutPath()
	out, oerr := cmd.CombinedOutput()

	if oerr != nil {
		log.Print("ERROR:")
		log.Panic(oerr)
	}

	postOut := ""
	if oerr == nil {
		postOut = s.getHookScriptResult("post_up")
	}

	// Yeah, there are better ways to do this.
	retval := preOut + "\n" + string(out) + "\n" + postOut
	return retval
}

func (s *Svn) SwitchTrunk() (string) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	preSw := s.getHookScriptResult("pre_sw")
	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/trunk")
	cmd.Dir = s.getCheckoutPath()
	out, oerr := cmd.CombinedOutput()

	if oerr != nil {
		log.Panic(oerr)
	}

	postSw := ""
	if oerr == nil {
		postSw = s.getHookScriptResult("post_sw")
	}

	retval := preSw + "\n" + string(out) + "\n" + postSw
	return retval
}

func (s *Svn) SwitchBranch(branchName string) (string) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	preSw := s.getHookScriptResult("pre_sw")
	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/branches/" + branchName)
	cmd.Dir = s.getCheckoutPath()

	out, oerr := cmd.CombinedOutput()

	if oerr != nil {
		log.Panic(oerr)
	}

	postSw := ""
	if oerr == nil {
		postSw = s.getHookScriptResult("post_sw")
	}

	retval := preSw + "\n" + string(out) + "\n" + postSw
	return retval
}

func (s *Svn) SwitchTag(tagName string) (string) {
	if s.shouldAlwaysRevert() {
		s.revertWorkingCopy()
	}

	preSw := s.getHookScriptResult("pre_sw")
	cmd := exec.Command(s.getBinPath(), "sw", s.getRepositoryPath() + "/tags/" + tagName)
	cmd.Dir = s.getCheckoutPath()
	out, oerr := cmd.CombinedOutput()

	if oerr != nil {
		log.Panic(oerr)
	}

	postSw := ""
	if oerr == nil {
		postSw = s.getHookScriptResult("post_sw")
	}

	retval := preSw + "\n" + string(out) + "\n" + postSw
	return retval
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
