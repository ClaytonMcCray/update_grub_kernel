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
	shell                 = flag.String("shell", "zsh", "Shell to execute in")
	overrideBackupFailure = flag.Bool("override_backup_failure", false, "Keep going if backing up grub_defaults fails")
	updateGrubPrg         = flag.String("update_grub_prg", "update-grub2", "Program to run to update grub after setting new default")
	runUpdateGrub         = flag.Bool("run_update_grub", true, "Whether or not to run update_grub_prg")
)

const (
	shellOpts      = "-c"
	grubDefaultKey = "GRUB_DEFAULT"
	searchPhrase   = "$menuentry_id_option "
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open(*grubCfg)
	if err != nil {
		log.Fatalf("could not open grub config: %s", err)
	}
	defer f.Close()

	grubCfgBody, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading grub config: %s", err)
	}

	kernels := findKernels(string(grubCfgBody))

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
		c := exec.Command(*shell, shellOpts, *updateGrubPrg)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	}
}

func findKernels(fileBody string) []string {
	lines := strings.Split(fileBody, "\n")
	kernels := []string{}

	for _, entry := range lines {
		if strings.Contains(entry, searchPhrase) {
			k, err := process(entry)
			if err != nil {
				log.Printf("error processing line %s. error: %s", entry, err)
			}

			kernels = append(kernels, strings.ReplaceAll(k, "'", "\""))
		}
	}

	return kernels
}

func backupDefaultsFile() error {
	f, err := os.Open(*grubDefaultsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	fbak, err := os.Create(*grubDefaultsFile + ".bak")
	defer fbak.Close()

	fbytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = fbak.Write(fbytes)
	return err
}

// readReducedDefaults reads the defaults file and removes blank lines.
func readReducedDefaults(f io.Reader) ([]string, error) {
	fb, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fb), "\n")
	reduced := []string{}
	for _, l := range lines {
		if l == "" || l == "\n" {
			continue
		}

		reduced = append(reduced, l)
	}

	return reduced, nil
}

func writeNewDefault(kernel string) error {
	f, err := os.OpenFile(*grubDefaultsFile, os.O_RDWR, 0000 /* not used, file exists */)
	if err != nil {
		log.Printf("opening %s: %s", *grubDefaultsFile, err)
		return err
	}
	defer f.Close()

	lines, err := readReducedDefaults(f)
	if err != nil {
		log.Printf("error reading %s: %s", f.Name(), err)
		return err
	}

	f.Truncate(0) // remove contents

	linesWithoutDefault := []string{}
	for _, l := range lines {
		if strings.HasPrefix(l, grubDefaultKey+"=") {
			continue
		}

		linesWithoutDefault = append(linesWithoutDefault, l)
	}

	linesWithoutDefault = append(linesWithoutDefault, grubDefaultKey+"="+kernel)

	_, err = f.WriteAt([]byte(strings.Join(linesWithoutDefault, "\n")+"\n"), 0)
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
	idx := strings.Index(line, searchPhrase)
	if idx < 0 {
		return "", fmt.Errorf("error cutFront: %s not found in %s", searchPhrase, line)
	}

	return line[idx+len(searchPhrase):], nil
}

func cutRear(line string) (string, error) {
	const end = " {"
	cut := strings.TrimSuffix(line, end)

	if cut == line {
		return "", fmt.Errorf("error cutRear: nothing was cut from %s", line)
	}

	return cut, nil
}
