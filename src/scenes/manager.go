package scenes

type SceneManager interface {
	GoTo(scene Scene)
	IsFadeInFinished() bool
	IsFading() bool
}
