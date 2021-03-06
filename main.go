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

const (
	shellOpts      = "-c"
	grubDefaultKey = "GRUB_DEFAULT"
	searchPhrase   = "$menuentry_id_option "
)

func main() {
	if err := Run(os.Stdin, os.Stderr, os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(stdin io.Reader, stderr, stdout io.Writer, args []string) error {
	log.SetOutput(stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		grubCfg               = flags.String("grub_cfg", "/boot/grub/grub.cfg", "Grub config file to parse")
		grubDefaultsFile      = flags.String("grub_defaults", "/etc/default/grub", "Where to write selected kernel")
		shell                 = flags.String("shell", "zsh", "Shell to execute in")
		overrideBackupFailure = flags.Bool("override_backup_failure", false, "Keep going if backing up grub_defaults fails")
		updateGrubPrg         = flags.String("update_grub_prg", "update-grub2", "Program to run to update grub after setting new default")
		runUpdateGrub         = flags.Bool("run_update_grub", true, "Whether or not to run update_grub_prg")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	grubCfgBody, err := readFile(*grubCfg)
	if err != nil {
		return fmt.Errorf("error opening grubCfg: %s", err)
	}

	kernels := findKernels(string(grubCfgBody))

	kernelSelection := kernels[userSelectsKernel(kernels)]

	if err := backupDefaultsFile(*grubDefaultsFile); err != nil {
		if *overrideBackupFailure {
			log.Printf("failure during backup, continuing: %s", err)
		} else {
			return fmt.Errorf("failure during backup: %s", err)
		}
	}

	if err := writeNewDefault(kernelSelection, *grubDefaultsFile); err != nil {
		return fmt.Errorf("error writing new default; note user must be root: %s", err)
	}

	if *runUpdateGrub {
		c := exec.Command(*shell, shellOpts, *updateGrubPrg)
		c.Stdout = stdout
		c.Stderr = stderr
		c.Run()
	}

	return nil
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not open grub config: %s", err)
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading grub config: %s", err)
	}

	return body, nil
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

func backupDefaultsFile(name string) error {
	contents, err := readFile(name)
	if err != nil {
		return err
	}

	fbak, err := os.Create(name + ".bak")
	defer fbak.Close()

	_, err = fbak.Write(contents)
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

func writeNewDefault(kernel string, defaultsName string) error {
	f, err := os.OpenFile(defaultsName, os.O_RDWR, 0000 /* not used, file exists */)
	if err != nil {
		log.Printf("opening %s: %s", defaultsName, err)
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
