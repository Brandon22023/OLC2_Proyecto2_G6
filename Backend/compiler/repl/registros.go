package repl

// Registros de propósito general
const (
    X0  = "x0"
    X1  = "x1"
    X2  = "x2"
    X3  = "x3"
    X4  = "x4"
    X5  = "x5"
    X6  = "x6"
    X7  = "x7"
    X8  = "x8"
    X9  = "x9"
    HP  = "x10"
    X11 = "x11"
    X12 = "x12"
    X13 = "x13"
    X14 = "x14"
    X15 = "x15"
    X16 = "x16"
    X17 = "x17"
    X18 = "x18"
    X19 = "x19"
    X20 = "x20"
    X21 = "x21"
    X22 = "x22"
    X23 = "x23"
    X24 = "x24"
    X25 = "x25"
    X26 = "x26"
    X27 = "x27"
    X28 = "x28"
    X29 = "x29"
    X30 = "x30"
)

// Registros de propósito especial
const (
    SP  = "sp"
    PC  = "pc"
    XZR = "xzr" // cero

    // Alias
    FP = "x29" // frame pointer (alias de x29)
    LR = "x30" // link register (alias de x30)
)

// Registros de punto flotante
const (
    D0  = "d0"
    D1  = "d1"
    D2  = "d2"
    D3  = "d3"
    D4  = "d4"
    D5  = "d5"
    D6  = "d6"
    D7  = "d7"
    D8  = "d8"
    D9  = "d9"
    D10 = "d10"
    D11 = "d11"
    D12 = "d12"
    D13 = "d13"
    D14 = "d14"
    D15 = "d15"
    D16 = "d16"
    D17 = "d17"
    D18 = "d18"
    D19 = "d19"
    D20 = "d20"
    D21 = "d21"
    D22 = "d22"
    D23 = "d23"
    D24 = "d24"
    D25 = "d25"
    D26 = "d26"
    D27 = "d27"
    D28 = "d28"
    D29 = "d29"
    D30 = "d30"
    D31 = "d31"
)

// Para strings (registros W)
const (
    W0 = "w0"
    W1 = "w1"
    W2 = "w2"
    W3 = "w3"
    W4 = "w4"
    W5 = "w5"
    W6 = "w6"
    W7 = "w7"
)