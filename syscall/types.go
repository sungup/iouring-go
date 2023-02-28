//go:build linux
// +build linux

package iouring_syscall

const (
	SYS_IO_URING_SETUP    = 425
	SYS_IO_URING_ENTER    = 426
	SYS_IO_URING_REGISTER = 427
)

// IORING Offset
const (
	IORING_OFF_SQ_RING uint64 = 0
	IORING_OFF_CQ_RING uint64 = 0x8000000
	IORING_OFF_SQES    uint64 = 0x10000000
)

const (
	IORING_OP_NOP uint8 = iota
	IORING_OP_READV
	IORING_OP_WRITEV
	IORING_OP_FSYNC
	IORING_OP_READ_FIXED
	IORING_OP_WRITE_FIXED
	IORING_OP_POLL_ADD
	IORING_OP_POLL_REMOVE
	IORING_OP_SYNC_FILE_RANGE
	IORING_OP_SENDMSG
	IORING_OP_RECVMSG
	IORING_OP_TIMEOUT
	IORING_OP_TIMEOUT_REMOVE
	IORING_OP_ACCEPT
	IORING_OP_ASYNC_CANCEL
	IORING_OP_LINK_TIMEOUT
	IORING_OP_CONNECT
	IORING_OP_FALLOCATE
	IORING_OP_OPENAT
	IORING_OP_CLOSE
	IORING_OP_FILES_UPDATE
	IORING_OP_STATX
	IORING_OP_READ
	IORING_OP_WRITE
	IORING_OP_FADVISE
	IORING_OP_MADVISE
	IORING_OP_SEND
	IORING_OP_RECV
	IORING_OP_OPENAT2
	IORING_OP_EPOLL_CTL
	IORING_OP_SPLICE
	IORING_OP_PROVIDE_BUFFERS
	IORING_OP_REMOVE_BUFFERS
	IORING_OP_TEE
	IORING_OP_SHUTDOWN
	IORING_OP_RENAMEAT
	IORING_OP_UNLINKAT
	IORING_OP_MKDIRAT
	IORING_OP_SYMLINKAT
	IORING_OP_LINKAT
	IORING_OP_MSG_RING
	IORING_OP_FSETXATTR
	IORING_OP_SETXATTR
	IORING_OP_FGETXATTR
	IORING_OP_GETXATTR
	IORING_OP_SOCKET
	IORING_OP_URING_CMD
	IORING_OP_SEND_ZC
	IORING_OP_SENDMSG_ZC

	/* this goes last, obviously */
	IORING_OP_LAST
)

const (
	IORING_SQ_NEED_WAKEUP uint32 = 1 << iota
	IORING_SQ_CQ_OVERFLOW
)

const (
	IOSQE_FLAGS_FIXED_FILE uint8 = 1 << iota
	IOSQE_FLAGS_IO_DRAIN
	IOSQE_FLAGS_IO_LINK
	IOSQE_FLAGS_IO_HARDLINK
	IOSQE_FLAGS_ASYNC
	IOSQE_FLAGS_BUFFER_SELECT
)

const IOSQE_SYNC_DATASYNC uint = 1
const IOSQE_TIMEOUT_ABS uint = 1
const IOSQE_SPLICE_F_FD_IN_FIXED = 1 << 31

type SubmissionQueueEntry interface {
	Opcode() uint8
	Reset()
	PrepOperation(op uint8, fd int32, addrOrSpliceOffIn uint64, len uint32, offsetOrCmdOp uint64)
	Fd() int32
	SetFdIndex(index int32)
	SetOpFlags(opflags uint32)
	SetUserData(userData uint64)
	SetFlags(flag uint8)
	CleanFlags(flags uint8)
	SetIoprio(ioprio uint16)
	SetBufIndex(bufIndex uint16)
	SetBufGroup(bufGroup uint16)
	SetPersonality(personality uint16)
	SetSpliceFdIn(fdIn int32)
}

type SQECore struct {
	opcode   uint8
	flags    uint8
	ioprio   uint16
	fd       int32
	offset   uint64
	addr     uint64
	len      uint32
	opFlags  uint32
	userdata uint64

	bufIndexOrGroup uint16
	personality     uint16
	spliceFdIn      int32
	extra           [2]uint64
}

func (sqe *SQECore) Opcode() uint8 {
	return sqe.opcode
}

func (sqe *SQECore) Reset() {
	*sqe = SQECore{}
}

func (sqe *SQECore) PrepOperation(op uint8, fd int32, addrOrSpliceOffIn uint64, len uint32, offsetOrCmdOp uint64) {
	sqe.opcode = op
	sqe.fd = fd
	sqe.addr = addrOrSpliceOffIn
	sqe.len = len
	sqe.offset = offsetOrCmdOp
}

func (sqe *SQECore) Fd() int32 {
	return sqe.fd
}

func (sqe *SQECore) SetFdIndex(index int32) {
	sqe.fd = index
	sqe.flags |= IOSQE_FLAGS_FIXED_FILE
}

func (sqe *SQECore) SetOpFlags(opflags uint32) {
	sqe.opFlags = opflags
}

func (sqe *SQECore) SetUserData(userData uint64) {
	sqe.userdata = userData
}

func (sqe *SQECore) SetFlags(flags uint8) {
	sqe.flags |= flags
}

func (sqe *SQECore) CleanFlags(flags uint8) {
	sqe.flags ^= flags
}

func (sqe *SQECore) SetIoprio(ioprio uint16) {
	sqe.ioprio = ioprio
}

func (sqe *SQECore) SetBufIndex(bufIndex uint16) {
	sqe.bufIndexOrGroup = bufIndex
}

func (sqe *SQECore) SetBufGroup(bufGroup uint16) {
	sqe.bufIndexOrGroup = bufGroup
}

func (sqe *SQECore) SetPersonality(personality uint16) {
	sqe.personality = personality
}

func (sqe *SQECore) SetSpliceFdIn(fdIn int32) {
	sqe.spliceFdIn = fdIn
}

type CompletionQueueEvent interface {
	UserData() uint64
	Result() int32
	Flags() uint32
	Clone() CompletionQueueEvent
}

type CQECore struct {
	userData uint64
	result   int32
	flags    uint32
}

func (cqe *CQECore) copyTo(dest *CQECore) {
	*dest = *cqe
}

func (cqe *CQECore) UserData() uint64 {
	return cqe.userData
}

func (cqe *CQECore) Result() int32 {
	return cqe.result
}

func (cqe *CQECore) Flags() uint32 {
	return cqe.flags
}

func (cqe *CQECore) Clone() CompletionQueueEvent {
	dest := &CQECore{}
	cqe.copyTo(dest)
	return dest
}

const IORING_FSYNC_DATASYNC uint32 = 1
const IORING_TIMEOUT_ABS uint32 = 1
