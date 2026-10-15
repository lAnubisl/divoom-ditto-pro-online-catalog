package playback

type Runner interface {
	Changed(immediate, takeover bool)
}
