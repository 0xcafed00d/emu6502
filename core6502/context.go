package core6502

type CPUMemory interface {
	Peek(addr uint16) uint8
	Poke(addr uint16, val uint8)
	PeekWord(addr uint16) uint16
	PokeWord(addr uint16, val uint16)
}

type CPURegisters interface {
	Flag(mask uint8) bool
	Flags() uint8
	RegA() uint8
	RegX() uint8
	RegY() uint8
	RegSP() uint8
	RegPC() uint16
	SetFlag(mask uint8, val bool)
	SetFlags(val uint8)
	SetRegA(val uint8)
	SetRegX(val uint8)
	SetRegY(val uint8)
	SetRegSP(val uint8)
	SetRegPC(val uint16)
}

type CPUContext interface {
	CPUMemory
	CPURegisters
}

// performs hard reset, resets all ram & registers
// attempts to write new reset vector
// resets stack pointer, and reloads PC from reset vector
func HardResetCPU(ctx CPUContext, resetVector uint16) {
	for n := 0; n < 0x10000; n++ {
		ctx.Poke(uint16(n), 0)
	}
	ctx.SetFlags(0)
	ctx.SetRegA(0)
	ctx.SetRegX(0)
	ctx.SetRegY(0)
	ctx.SetRegSP(0xff)
	ctx.PokeWord(Vector_RST, resetVector)
	ctx.SetRegPC(ctx.PeekWord(Vector_RST))
}

// performs soft reset. resets stack pointer, and reloads PC from reset vector
func SoftResetCPU(ctx CPUContext) {
	ctx.SetRegSP(0xff)
	ctx.SetRegPC(ctx.PeekWord(Vector_RST))
}

func HiByte(val uint16) uint8 {
	return uint8(val >> 8)
}

func LoByte(val uint16) uint8 {
	return uint8(val)
}

func MakeWord(hi, lo uint8) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}

func Push8(ctx CPUContext, val uint8) {
	sp := ctx.RegSP()
	ctx.Poke(0x100+uint16(sp), val)
	ctx.SetRegSP(sp - 1)
}

func Push16(ctx CPUContext, val uint16) {
	Push8(ctx, uint8(val>>8))
	Push8(ctx, uint8(val))
}

func Pop8(ctx CPUContext) uint8 {
	sp := ctx.RegSP() + 1
	val := ctx.Peek(0x100 + uint16(sp))
	ctx.SetRegSP(sp)
	return val
}

func Pop16(ctx CPUContext) uint16 {
	var val uint16 = uint16(Pop8(ctx))
	val |= uint16(Pop8(ctx)) << 8
	return val
}

// CPU context. Conains CPU registers and 64K Ram
type BasicCPUContext struct {
	reg struct {
		a, x, y   uint8
		sp, flags uint8
		pc        uint16
	}

	ram [0x10000]uint8
}

// Masks for flag register
const (
	Flag_C      uint8 = 1 << iota // Carry
	Flag_Z      uint8 = 1 << iota // Zero
	Flag_I      uint8 = 1 << iota // IRQ Disable
	Flag_D      uint8 = 1 << iota // Decimal Mode
	Flag_B      uint8 = 1 << iota // Break Command
	Flag_unused uint8 = 1 << iota // not used
	Flag_V      uint8 = 1 << iota // Overflow
	Flag_N      uint8 = 1 << iota // Negative
)

const (
	Vector_NMI uint16 = 0xfffa
	Vector_RST uint16 = 0xfffc
	Vector_IRQ uint16 = 0xfffe
)

func (c *BasicCPUContext) Flag(mask uint8) bool {
	return (c.reg.flags & mask) != 0
}

func (c *BasicCPUContext) Flags() uint8 {
	return c.reg.flags
}

func (c *BasicCPUContext) RegA() uint8 {
	return c.reg.a
}

func (c *BasicCPUContext) RegX() uint8 {
	return c.reg.x
}

func (c *BasicCPUContext) RegY() uint8 {
	return c.reg.y
}

func (c *BasicCPUContext) RegSP() uint8 {
	return c.reg.sp
}

func (c *BasicCPUContext) RegPC() uint16 {
	return c.reg.pc
}

func (c *BasicCPUContext) SetFlag(mask uint8, val bool) {
	if val {
		c.reg.flags |= mask
	} else {
		c.reg.flags &^= mask
	}
}

func (c *BasicCPUContext) SetFlags(val uint8) {
	c.reg.flags = val
}

func (c *BasicCPUContext) SetRegA(val uint8) {
	c.reg.a = val
}

func (c *BasicCPUContext) SetRegX(val uint8) {
	c.reg.x = val
}

func (c *BasicCPUContext) SetRegY(val uint8) {
	c.reg.y = val
}

func (c *BasicCPUContext) SetRegSP(val uint8) {
	c.reg.sp = val
}

func (c *BasicCPUContext) SetRegPC(val uint16) {
	c.reg.pc = val
}

func (c *BasicCPUContext) Peek(addr uint16) uint8 {
	return c.ram[addr]
}

func (c *BasicCPUContext) Poke(addr uint16, val uint8) {
	c.ram[addr] = val
}

func (c *BasicCPUContext) PeekWord(addr uint16) uint16 {
	var val uint16 = uint16(c.Peek(addr+1)) << 8
	val |= uint16(c.Peek(addr))
	return val
}

func (c *BasicCPUContext) PokeWord(addr uint16, val uint16) {
	c.Poke(addr, uint8(val))
	c.Poke(addr+1, uint8(val>>8))
}
