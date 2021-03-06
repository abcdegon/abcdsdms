/******************************************************************************/
/* dms.go                                                                     */
/*                                                                            */
/* Usage: dms -f filing directory -i JSON index file -l log file \            */
/*              -s storage directory                                          */
/*                                                                            */
/*           -f : directory (location) of the files before sorting            */
/*           -i : JSON index file of the filing system                        */
/*           -l : log file of the program                                     */
/*           -s : directory in which files are stored sorted                  */
/*                                                                            */
/* Purpose: Filing system for files to be saved according to a predefined     */
/*          naming scheme. Naming Scheme is described in the wiki of the      */
/*          github project page.                                              */
/*                                                                            */
/* Author: abcddev@yahoo.com                                                  */
/*                                                                            */
/* Version: 0.0.0-4                                                           */
/*                                                                            */
/* Revision: 0 - initialize the source code files.                            */
/*           1 - added comments and finished checkCmdLineArgs()               */
/*           2 - edited checkCmdLineArgs to adopt logging format              */
/*           3 - added comments and finished dmsLog()                         */
/*           4 - changed dmsLog() to be more flexible                         */
/******************************************************************************/
package main

// imports
import (
	"errors"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const ()

/******************************************************************************/
/* TODO:                                                                      */
/* This is the todo list of this project.                                     */
/*                                                                            */
/* [X] create a log format                                                    */
/*     [X] adopt logging in checkCmdLineArgs()                                */
/*     [X] adopt logging in main()                                            */
/*     [ ] create a format for the transaction id                             */
/* [ ] read filing directory                                                  */
/******************************************************************************/

func checkCmdLineArgs(args []string) (fi, in, lf, st string, err error) {
	/**************************************************************************/
	/* Description of return values                                           */
	/* fi string - location of files before sorting                           */
	/* in string - json index file                                            */
	/* lf string - name of the log file                                       */
	/* st string - location of files when sorted                              */
	/* err error - error if any                                               */
	/**************************************************************************/

	/**************************************************************************/
	/* Local variables                                                        */
	/*        i - int              - counter                                  */
	/**************************************************************************/
	var i int

	/**************************************************************************/
	/* Check every single commandline argument. If a valid option is found    */
	/* set the according value. Invalid options are ignored.                  */
	/**************************************************************************/
	for i = 0; i < len(args); i++ {
		switch args[i] {
		case "-f":
			fi = args[i+1]
		case "-i":
			in = args[i+1]
		case "-l":
			lf = args[i+1]
		case "-s":
			st = args[i+1]
		}
	}

	// Check if all the necessary parameters are set or set error code.
	if fi == "" {
		err = errors.New("100001")
	} else {
		if in == "" {
			err = errors.New("100002")
		} else {
			if lf == "" {
				err = errors.New("100003")
			} else {
				if st == "" {
					err = errors.New("100004")
				}
			}
		}
	}

	return
}

func dmsLog(file *os.File, id int, tid string, src string) {
	/*************************************************************************/
	/* Local variables                                                       */
	/*      msgs - map[string]string - map of all the log messages           */
	/*    fields - log.Fields        - fields to be logged                   */
	/*        sl - log.Level         - the syslog level                      */
	/*************************************************************************/
	var msgs = map[int]string{
		100001: "Missing command line argument -f",
		100002: "Missing command line argument -i",
		100003: "Missing command line argument -l",
		100004: "Missing command line argument -s",
		200001: "Error with your logfile, set file to STDOUT",
		700001: "Abcdsdms started"}
	var fields = log.Fields{}
	var sl = log.FatalLevel
	var level = id / 100000

	/**************************************************************************/
	/* Set the fields which will be logged. Their values are given by the     */
	/* function parameters.                                                   */
	/**************************************************************************/
	fields["id"] = fmt.Sprintf("%s-%d-%d", "AbcdsDMS", level, id)
	fields["transaction"] = tid
	fields["src"] = src

	/**************************************************************************/
	/* This switch command is used to set the syslog level of the upcoming    */
	/* message. It is important because that's how the syslog level is logged */
	/* itself.                                                                */
	/**************************************************************************/
	switch level {
	case 1:
		sl = log.PanicLevel
	case 2:
		sl = log.FatalLevel
	case 3:
		sl = log.ErrorLevel
	case 4:
		sl = log.WarnLevel
	case 5, 6:
		sl = log.InfoLevel
	case 7:
		sl = log.DebugLevel
	}

	// If there is no log file opened I will send messages to STDOUT.
	if file == nil {
		file = os.Stdout
	}

	/**************************************************************************/
	/* The following command generate the log message that will be written    */
	/* either to STDOUT or to file. The syslog will be written in JSON which  */
	/* to use log analysis tools like Splunk>, Graylog, etc.                  */
	/**************************************************************************/
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(sl)
	log.WithFields(fields).Log(sl, msgs[id])
}

func main() {
	/**************************************************************************/
	/* Descriptiion of local variables                                        */
	/*    filing : the directory where files are stored before dms runs       */
	/*     index : the json index of abcdsdms                                 */
	/*     logfn : the filename of the logfile                                */
	/*   storage : the directory where files are stored after dms ran         */
	/*       err : any error returned by any function                         */
	/*     *logf : a pointer to the logfile handle                            */
	/**************************************************************************/
	var filing, index, logfn, storage string
	var err error
	var logf *os.File
	var id int

	/**************************************************************************/
	/* Parse the commandline arguments. If any error occurs write a syslog    */
	/* message of level fatal and exit the program.                           */
	/**************************************************************************/
	filing, index, logfn, storage, err = checkCmdLineArgs(os.Args)

	if err != nil {
		id, _ = strconv.Atoi(err.Error())
		dmsLog(os.Stdout, id, "0", "main")
	}

	/**************************************************************************/
	/* Open the log file and prepare it for logging. Create it if it not      */
	/* exists or append if it exists or open it for writing only. On creation */
	/* set it mode to 0644. If any error occurs write a syslog message of     */
	/* level fatal and exit the program.                                      */
	/**************************************************************************/
	logf, err = os.OpenFile(logfn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		dmsLog(os.Stdout, 200001, "0", "main")
		logf = os.Stdout
	}

	dmsLog(logf, 700001, "0", "main")
	fmt.Println(filing, index, logfn, storage)
	// get a list of all files in -f

	/* Nassi Shneiderman Diagram
		 **************************************
		 *****  Commandline Arguments ok? *****
		 ***         true           | false ***
		 ***************************|**********
		 ** start logging subsystem | end    **
		 ***************************|        **
	     ** list all files in -f    |        **
		 ***************************|        **
		 ** For all files           |        **
		 **   **********************|        **
		 **   **   filename ok?   **|        **
		 **   **  true  |  false  **|        **
		 **   **********|***********|        **
		 **   ** md5    | error   **|        **
		 **   **********|         **|        **
		 **   ** save   |         **|        **
		 **   ** to -s  |         **|        **
		 **   **********|         **|        **
		 **   ** write  |         **|        **
		 **   ** index  |         **|        **
		 **   **********|***********|        **
		 ***************************|**********
		 ** End of Programm                  **
		 **************************************
	*/
}
