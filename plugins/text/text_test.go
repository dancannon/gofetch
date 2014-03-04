package text

import (
	"github.com/davecgh/go-spew/spew"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClassifyBlock(t *testing.T) {
	Convey("Subject: Classify block", t, func() {
		Convey("Given some input which is mostly unlinked text", func() {
			block := TextBlock{}
			block.AddText("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a", false)
			block.AddText("galley", true)
			block.AddText(" of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of", false)
			block.AddText("Lorem Ipsum", true)
			block.AddText(".", false)
			block.Flush()

			Convey("When the block is flushed", func() {
				block.Classify()

				spew.Dump(block)
				Convey("The resulting type is 'Content'", func() {
					So(block.Type, ShouldEqual, Content)
				})
			})
		})
	})
}
