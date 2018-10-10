package grafeas

import (
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
)

const noNewOcc = "no new occurrences generated"

// Item implements a MetadataItem.
type Item struct {
	Occurrence *grafeaspb.Occurrence // The Occurrence this Item wraps.
}

// Name returns the name of the group of Item.
func (item *Item) Name() string {
	return item.Occurrence.NoteName
}

// String returns a string version of this Item.
func (item *Item) String() string {
	if nil != item.Occurrence {
		return item.Occurrence.String()
	}

	return noNewOcc
}
