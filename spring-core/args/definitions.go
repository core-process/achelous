package args

import "time"

type ArgProgram int8

const (
	ArgProgramSendmail   ArgProgram = 1
	ArgProgramNewaliases ArgProgram = 2
	ArgProgramMailq      ArgProgram = 3
)

// See https://www.sendmail.org/~ca/email/man/sendmail.html for more.

// -Btype      Set the body type to type. Current legal values 7BIT or
//             8BITMIME.
type SmArg_B int8

const (
	SmArg_B_7Bit     SmArg_B = 0
	SmArg_B_8BitMime SmArg_B = 1
)

// -N dsn      Set delivery status notification conditions to dsn, which can
//             be `never' for no notifications or a comma separated list of
//             the values `failure' to be notified if delivery failed,
//             `delay' to be notified if delivery is delayed, and `success'
//             to be notified when the message is successfully delivered.
type SmArg_N int8

const (
	SmArg_N_Never   SmArg_N = 0
	SmArg_N_Failure SmArg_N = 1
	SmArg_N_Delay   SmArg_N = 2
	SmArg_N_Success SmArg_N = 4
)

// -pprotocol  Set the name of the protocol used to receive the message.
//             This can be a simple protocol name such as ``UUCP'' or a pro-
//             tocol and hostname, such as ``UUCP:ucbvax''.
type SmArg_p struct {
	Protocol string
	Hostname *string
}

// -R return   Set the amount of the message to be returned if the message
//             bounces. The return parameter can be `full' to return the
//             entire message or `hdrs' to return only the headers.
type SmArg_R int8

const (
	SmArg_R_Full SmArg_R = 0
	SmArg_R_Hdrs SmArg_R = 1
)

// -O option=value
//             Set option option to the specified value. This form uses long
//             names. See below for more details.
type SmArg_O struct {

	// AliasFile=file
	//             Use alternate alias file.
	Opt_AliasFile *string

	// HoldExpensive
	//             On mailers that are considered ``expensive'' to connect to,
	//             don't initiate immediate connection. This requires queueing.
	Opt_HoldExpensive bool

	// CheckpointInterval=N
	//             Checkpoint the queue file after every N successful deliveries
	//             (default 10). This avoids excessive duplicate deliveries
	//             when sending to long mailing lists interrupted by system
	//             crashes.
	Opt_CheckpointInterval *int16

	// DeliveryMode=x
	//             Set the delivery mode to x. Delivery modes are `i' for inter-
	//             active (synchronous) delivery, `b' for background (asyn-
	//             chronous) delivery, `q' for queue only - i.e., actual deliv-
	//             ery is done the next time the queue is run, and `d' for de-
	//             ferred - the same as `q' except that database lookups (no-
	//             tably DNS and NIS lookups) are avoided.
	// TODO: fix data type
	Opt_DeliveryMode *rune

	// ErrorMode=x
	//             Set error processing to mode x. Valid modes are `m' to mail
	//             back the error message, `w' to ``write'' back the error mes-
	//             sage (or mail it back if the sender is not logged in), `p' to
	//             print the errors on the terminal (default), `q' to throw away
	//             error messages (only exit status is returned), and `e' to do
	//             special processing for the BerkNet. If the text of the mes-
	//             sage is not mailed back by modes `m' or `w' and if the sender
	//             is local to this machine, a copy of the message is appended
	//             to the file dead.letter in the sender's home directory.
	// TODO: fix data type
	Opt_ErrorMode *rune

	// SaveFromLine
	//             Save UNIX-style From lines at the front of messages.
	Opt_SaveFromLine bool

	// MaxHopCount= N
	//             The maximum number of times a message is allowed to ``hop''
	//             before we decide it is in a loop.
	Opt_MaxHopCount *int16

	// IgnoreDots  Do not take dots on a line by themselves as a message termi-
	//             nator.
	Opt_IgnoreDots bool

	// SendMimeErrors
	//             Send error messages in MIME format. If not set, the DSN (De-
	//             livery Status Notification) SMTP extension is disabled.
	Opt_SendMimeErrors bool

	// ConnectionCacheTimeout=timeout
	//             Set connection cache timeout.
	Opt_ConnectionCacheTimeout *time.Duration

	// ConnectionCacheSize=N
	//             Set connection cache size.
	Opt_ConnectionCacheSize *int16

	// LogLevel=n  The log level.
	Opt_LogLevel *int16

	// MeToo       Send to ``me'' (the sender) also if I am in an alias expan-
	//             sion.
	Opt_MeToo bool

	// CheckAliases
	//             Validate the right hand side of aliases during a newalias-
	//             es(1) command.
	Opt_CheckAliases bool

	// OldStyleHeaders
	//             If set, this message may have old style headers. If not set,
	//             this message is guaranteed to have new style headers (i.e.,
	//             commas instead of spaces between addresses). If set, an
	//             adaptive algorithm is used that will correctly determine the
	//             header format in most cases.
	Opt_OldStyleHeaders bool

	// QueueDirectory=queuedir
	//             Select the directory in which to queue messages.
	Opt_QueueDirectory *string

	// StatusFile=file
	//             Save statistics in the named file.
	Opt_StatusFile *string

	// Timeout.queuereturn=time
	//             Set the timeout on undelivered messages in the queue to the
	//             specified time. After delivery has failed (e.g., because of
	//             a host being down) for this amount of time, failed messages
	//             will be returned to the sender. The default is five days.
	Opt_TimeoutQueueReturn *time.Duration

	// UserDatabaseSpec=userdatabase
	//             If set, a user database is consulted to get forwarding infor-
	//             mation. You can consider this an adjunct to the aliasing
	//             mechanism, except that the database is intended to be dis-
	//             tributed; aliases are local to a particular host. This may
	//             not be available if your sendmail does not have the USERDB
	//             option compiled in.
	Opt_UserDatabaseSpec *string

	// ForkEachJob
	//             Fork each job during queue runs. May be convenient on memo-
	//             ry-poor machines.
	Opt_ForkEachJob bool

	// SevenBitInput
	//             Strip incoming messages to seven bits.
	Opt_SevenBitInput bool

	// EightBitMode=mode
	//             Set the handling of eight bit input to seven bit destinations
	//             to mode: m (mimefy) will convert to seven-bit MIME format, p
	//             (pass) will pass it as eight bits (but violates protocols),
	//             and s (strict) will bounce the message.
	// TODO: fix data type
	Opt_EightBitMode *rune

	// MinQueueAge=timeout
	//             Sets how long a job must ferment in the queue between at-
	//             tempts to send it.
	Opt_MinQueueAge *time.Duration

	// DefaultCharSet=charset
	//             Sets the default character set used to label 8-bit data that
	//             is not otherwise labelled.
	Opt_DefaultCharSet *string

	// DialDelay=sleeptime
	//             If opening a connection fails, sleep for sleeptime seconds
	//             and try again. Useful on dial-on-demand sites.
	Opt_DialDelay *time.Duration

	// NoRecipientAction=action
	//             Set the behaviour when there are no recipient headers (To:,
	//             Cc: or Bcc:) in the message to action: none leaves the mes-
	//             sage unchanged, add-to adds a To: header with the envelope
	//             recipients, add-apparently-to adds an Apparently-To: header
	//             with the envelope recipients, add-bcc adds an empty Bcc:
	//             header, and add-to-undisclosed adds a header reading `To:
	//             undisclosed-recipients:;'.
	// TODO: fix data type
	Opt_NoRecipientAction *string

	// MaxDaemonChildren=N
	//             Sets the maximum number of children that an incoming SMTP
	//             daemon will allow to spawn at any time to N.
	Opt_MaxDaemonChildren *int16

	// ConnectionRateThrottle=N
	//             Sets the maximum number of connections per second to the SMTP
	//             port to N.
	Opt_ConnectionRateThrottle *int16
}

type SmArgs struct {

	// -Btype      Set the body type to type. Current legal values 7BIT or
	//             8BITMIME.
	Arg_B *SmArg_B `args:"value"`

	// -ba         Go into ARPANET mode. All input lines must end with a CR-LF,
	//             and all messages will be generated with a CR-LF at the end.
	//             Also, the ``From:'' and ``Sender:'' fields are examined for
	//             the name of the sender.
	Arg_ba bool `args:"flag"`

	// -bd         Run as a daemon. This requires Berkeley IPC. Sendmail will
	//             fork and run in background listening on socket 25 for incom-
	//             ing SMTP connections. This is normally run from /etc/rc.
	Arg_bd bool `args:"flag"`

	// -bD         Same as -bd except runs in foreground.
	Arg_bD bool `args:"flag"`

	// -bh         Print the persistent host status database.
	Arg_bh bool `args:"flag"`

	// -bH         Purge the persistent host status database.
	Arg_bH bool `args:"flag"`

	// -bi         Initialize the alias database.
	Arg_bi bool `args:"flag"`

	// -bm         Deliver mail in the usual way (default).
	Arg_bm bool `args:"flag"`

	// -bp         Print a listing of the queue.
	Arg_bp bool `args:"flag"`

	// -bs         Use the SMTP protocol as described in RFC821 on standard in-
	//             put and output. This flag implies all the operations of the
	//             -ba flag that are compatible with SMTP.
	Arg_bs bool `args:"flag"`

	// -bt         Run in address test mode. This mode reads addresses and
	//             shows the steps in parsing; it is used for debugging configu-
	//             ration tables.
	Arg_bt bool `args:"flag"`

	// -bv         Verify names only - do not try to collect or deliver a mes-
	//             sage. Verify mode is normally used for validating users or
	//             mailing lists.
	Arg_bv bool `args:"flag"`

	// -Cfile      Use alternate configuration file. Sendmail refuses to run as
	//             root if an alternate configuration file is specified.
	Arg_C *string `args:"value"`

	// -dX         Set debugging value to X.
	Arg_d *int16 `args:"value"`

	// -Ffullname  Set the full name of the sender.
	Arg_F *string `args:"value"`

	// -fname      Sets the name of the ``from'' person (i.e., the sender of the
	//             mail). -f can only be used by ``trusted'' users (normally
	//             root, daemon, and network) or if the person you are trying to
	//             become is the same as the person you are.
	Arg_f *string `args:"value"`

	// -hN         Set the hop count to N. The hop count is incremented every
	//             time the mail is processed. When it reaches a limit, the
	//             mail is returned with an error message, the victim of an
	//             aliasing loop. If not specified, ``Received:'' lines in the
	//             message are counted.
	Arg_h *int16 `args:"value"`

	// -i          Ignore dots alone on lines by themselves in incoming mes-
	//             sages. This should be set if you are reading data from a
	//             file.
	Arg_i bool `args:"flag"`

	// -N dsn      Set delivery status notification conditions to dsn, which can
	//             be `never' for no notifications or a comma separated list of
	//             the values `failure' to be notified if delivery failed,
	//             `delay' to be notified if delivery is delayed, and `success'
	//             to be notified when the message is successfully delivered.
	Arg_N *SmArg_N `args:"value"`

	// -n          Don't do aliasing.
	Arg_n bool `args:"flag"`

	// -O option=value
	//             Set option option to the specified value. This form uses long
	//             names. See below for more details.
	Arg_O SmArg_O `args:"value"`

	// -ox value   Set option x to the specified value. This form uses single
	//             character names only. The short names are not described in
	//             this manual page; see the Sendmail Installation and Operation
	//             Guide for details.
	// NOTE: we support -oi only for now
	Arg_oi bool `args:"flag"`

	// -pprotocol  Set the name of the protocol used to receive the message.
	//             This can be a simple protocol name such as ``UUCP'' or a pro-
	//             tocol and hostname, such as ``UUCP:ucbvax''.
	Arg_p *SmArg_p `args:"value"`

	// -q[time]    Processed saved messages in the queue at given intervals. If
	//             time is omitted, process the queue once. Time is given as a
	//             tagged number, with `s' being seconds, `m' being minutes, `h'
	//             being hours, `d' being days, and `w' being weeks. For exam-
	//             ple, `-q1h30m' or `-q90m' would both set the timeout to one
	//             hour thirty minutes. If time is specified, sendmail will run
	//             in background. This option can be used safely with -bd.
	// NOTE: currently we support -q as a flag only!
	Arg_q bool `args:"flag"`

	// -qIsubstr   Limit processed jobs to those containing substr as a sub-
	//             string of the queue id.
	Arg_qI *string `args:"value"`

	// -qRsubstr   Limit processed jobs to those containing substr as a sub-
	//             string of one of the recipients.
	Arg_qR *string `args:"value"`

	// -qSsubstr   Limit processed jobs to those containing substr as a sub-
	//             string of the sender.
	Arg_qS *string `args:"value"`

	// -R return   Set the amount of the message to be returned if the message
	//             bounces. The return parameter can be `full' to return the
	//             entire message or `hdrs' to return only the headers.
	Arg_R *SmArg_R `args:"value"`

	// -rname      An alternate and obsolete form of the -f flag.
	Arg_r *string `args:"value"`

	// -t          Read message for recipients. To:, Cc:, and Bcc: lines will
	//             be scanned for recipient addresses. The Bcc: line will be
	//             deleted before transmission.
	Arg_t bool `args:"flag"`

	// -U          Initial (user) submission. This should always be set when
	//             called from a user agent such as Mail or exmh and never be
	//             set when called by a network delivery agent such as rmail.
	Arg_U bool `args:"flag"`

	// -V envid    Set the original envelope id. This is propagated across SMTP
	//             to servers that support DSNs and is returned in DSN-compliant
	//             error messages.
	Arg_V *string `args:"value"`

	// -v          Go into verbose mode. Alias expansions will be announced,
	//             etc.
	Arg_v bool `args:"flag"`

	// -X logfile  Log all traffic in and out of mailers in the indicated log
	//             file. This should only be used as a last resort for debug-
	//             ging mailer bugs. It will log a lot of data very quickly.
	Arg_X *string `args:"value"`
}

type MqArgs struct {

	// -v          Go into verbose mode. Alias expansions will be announced,
	//             etc.
	Arg_v bool `args:"flag"`
}
