package getprofile

import (
    "regexp"
    "log"
    "os/exec"
    "path"
    "path/filepath"
    "os"
    "errors"
)

type gitSyncer struct{}

var gitSyncerInstance *gitSyncer = &gitSyncer{}

func (sync *gitSyncer) Supports(repoUrl string) bool {
    if matched, err := regexp.MatchString(".+@.+:.+", repoUrl); err != nil {
        log.Printf("Error checking gitSyncer.Supports(%v)", repoUrl)
        return false
    } else {
        return matched
    }
}

func (sync *gitSyncer) Init() error {
    if repoUrl, err := getConfig(); err != nil {
        return err
    } else if _, err := os.Stat(path.Join(repoPath, ".git")); os.IsNotExist(err) {
        return exec.Command("git", "clone", repoUrl, path.Join(basePath, "repo")).Run()
    } else {
        return nil // Already cloned
    }
}

func (sync *gitSyncer) Track(absPath string, relPath string) error {
    dest := path.Join(repoPath, relPath)
    if err := os.MkdirAll(path.Dir(dest), 0700); err != nil {
        return err
    }

    if err := exec.Command("cp", absPath, dest).Run(); err != nil {
        return err
    }

    cmd := exec.Command("git", "add", relPath)
    cmd.Dir = repoPath
    return cmd.Run()
}

func (sync *gitSyncer) Untrack(relPath string) error {
    cmd := exec.Command("git", "rm", "-f", relPath)
    cmd.Dir = repoPath
    return cmd.Run()
}

func (sync *gitSyncer) Out() error {
    if err := filepath.Walk(repoPath, makeWalkFunc(copyToRepo)); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")

        cmd := exec.Command("git", "diff", "HEAD", "--quiet")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        } else if cmd.ProcessState.Success() {
            return nil // No changes
        }

        cmd = exec.Command("git", "add", "-u")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        }

        cmd = exec.Command("git", "commit", "-m", "auto-commit")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        }

        cmd = exec.Command("git", "push")
        cmd.Dir = repoPath
        return cmd.Run()
    }
}

func makeWalkFunc(processor func(homePath string, absPath string, relPath string) error) filepath.WalkFunc {
    return func(absPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        } else if info.IsDir() {
            return nil
        } else if relPath, err := filepath.Rel(repoPath, absPath); err != nil {
            return err
        } else if !filepath.HasPrefix(relPath, ".git/") {
            return processor(absPath, relPath)
        } else {
            return nil
        }
    }
}

func (sync *gitSyncer) In() error {
    prevSha, err := gitSha(repoPath)
    if err != nil {
        return err
    }

    cmd := exec.Command("git", "pull")
    cmd.Dir = repoPath
    if err := cmd.Run(); err != nil {
        return err
    }

    nowSha, err := gitSha(repoPath)
    if err != nil {
        return err
    }

    if prevSha == nowSha {
        return nil
    }

    return filepath.Walk(path.Join(basePath, "repo"), copyToLocal)
}

func gitSha(repoPath string) (string, error) {
    cmd := exec.Command("git", "rev-parse", "HEAD")
    cmd.Dir = repoPath
    if err := cmd.Run(); err != nil {
        return "", err
    }
    bytes, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func copyToRepo(absPath string, relPath string) error {
    return errors.New("TODO") // TODO
}

func copyToLocal(absPath string, relPath string) error {
    return errors.New("TODO") // TODO
}

