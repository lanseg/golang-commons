package common

func IfNull[T any](a *T, b *T) *T {
    if a == nil {
        return b
    }
    return a
}

func IfEmpty(a string, b string) string {
    if a == "" {
        return b
    }
    return a
}
