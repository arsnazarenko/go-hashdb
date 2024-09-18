package util



func Assert(ok bool, msg string) {
    if !ok {
        panic(msg)
    }
}

