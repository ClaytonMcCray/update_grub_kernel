package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	grubCfg               = flag.String("grub_cfg", "/boot/grub/grub.cfg", "Grub config file to parse")
	grubDefaultsFile      = flag.String("grub_defaults", "/etc/default/grub", "Where to write selected kernel")
	grep                  = flag.String("grep_prg", "rg", "Grep-like program")
	shell                 = flag.String("shell", "zsh", "Shell to execute in")
	overrideBackupFailure = flag.Bool("override_backup_failure", false, "Keep going if backing up grub_defaults fails")
	updateGrubPrg         = flag.String("update_grub_prg", "update-grub2", "Program to run to update grub after setting new default")
	runUpdateGrub         = flag.Bool("run_update_grub", true, "Whether or not to run update_grub_prg")
)

const (
	searchOpts     = "-F"
	shellOpts      = "-c"
	grubDefaultKey = "GRUB_DEFAULT"
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := exec.Command(*shell, shellOpts, fmt.Sprintf("%s %s %s %s", *grep, searchOpts, searchPhraseForGrep(), *grubCfg))
	e, err := c.Output()
	if err != nil {
		log.Printf("Output: %s", e)
		log.Printf("Command run: %s", c.String())
		log.Fatal(err)
	}

	entries := strings.Split(string(e), "\n")
	kernels := []string{}
	for _, entry := range entries {
		processed, err := process(entry)
		if err != nil {
			log.Printf("error processing line %s, skipping", entry)
		}

		kernels = append(kernels, processed)
	}

	kernelSelection := kernels[userSelectsKernel(kernels)]

	if err = backupDefaultsFile(); err != nil {
		if *overrideBackupFailure {
			log.Printf("failure during backup, continuing: %s", err)
		} else {
			log.Fatalf("failure during backup: %s", err)
		}
	}

	if err = writeNewDefault(kernelSelection); err != nil {
		log.Fatalf("error writing new default; note user must be root: %s", err)
	}

	if *runUpdateGrub {
		c = exec.Command(*shell, shellOpts, *updateGrubPrg)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	}
}

func backupDefaultsFile() error {
	f, err := os.Open(*grubDefaultsFile)
	defer f.Close()
	if err != nil {
		return err
	}

	fbak, err := os.Create(*grubDefaultsFile + ".bak")
	defer fbak.Close()

	fbytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = fbak.Write(fbytes)
	return err
}

func writeNewDefault(kernel string) error {
	f, err := os.Open(*grubDefaultsFile)
	defer f.Close()
	if err != nil {
		log.Printf("opening %s: %s", *grubDefaultsFile, err)
		return err
	}

	inf, err := f.Stat()
	if err != nil {
		log.Printf("error getting stat for %s: %s", f.Name(), err)
	}

	permForDefaults := inf.Mode()

	fb, err := io.ReadAll(f)
	if err != nil {
		log.Printf("reading from %s: %s", f.Name(), err)
		return err
	}

	lines := strings.Split(string(fb), "\n")
	linesWithoutDefault := []string{}

	for _, l := range lines {
		if strings.HasPrefix(l, grubDefaultKey+"=") {
			continue
		}

		linesWithoutDefault = append(linesWithoutDefault, l)
	}

	f.Close() // close from when it was open to read

	linesWithoutDefault = append(linesWithoutDefault, grubDefaultKey+"="+strings.ReplaceAll(kernel, "'", "\""))
	f, err = os.OpenFile(*grubDefaultsFile, os.O_RDWR, permForDefaults)
	defer f.Close()
	if err != nil {
		log.Printf("opening %s: %s", *grubDefaultsFile, err)
		return err
	}

	_, err = f.Write([]byte(strings.Join(linesWithoutDefault, "\n") + "\n"))
	if err != nil {
		log.Printf("error writing %s: %s", f.Name(), err)
		return err
	}

	return nil
}

func userSelectsKernel(kernels []string) int {
	fmt.Println("Select the index of the kernel you wish to make default:")
	for i, e := range kernels {
		fmt.Printf("%d. %s\n", i, e)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Index: ")
		iStr, _ := reader.ReadString('\n')
		idx, err := strconv.Atoi(strings.Trim(iStr, "\n"))

		if err == nil {
			return idx
		}

		log.Printf("Got error converting to int, try again: %s", err)
	}
}

func searchPhraseGolang() string {
	return "$menuentry_id_option "
}

func searchPhraseForGrep() string {
	return fmt.Sprintf("\"\\%s\"", searchPhraseGolang())
}

func process(line string) (string, error) {
	tmp, err := cutFront(line)
	if err != nil {
		return "", err
	}

	tmp, err = cutRear(tmp)
	if err != nil {
		return "", err
	}

	return tmp, nil
}

func cutFront(line string) (string, error) {
	idx := strings.Index(line, searchPhraseGolang())
	if idx < 0 {
		return "", fmt.Errorf("error cutFront: %s not found in %s", searchPhraseGolang(), line)
	}

	return line[idx+len(searchPhraseGolang()):], nil
}

func cutRear(line string) (string, error) {
	const end = " {"
	cut := strings.TrimSuffix(line, end)

	if cut == line {
		return "", fmt.Errorf("error cutRear: nothing was cut from %s", line)
	}

	return cut, nil
}
