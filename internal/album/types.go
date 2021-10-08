package album

import "fmt"

// Album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Albums slice to seed record album data.
var Albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// Handles Add ALbum
func AddAlbum(a Album) {
	// Add the new album to the slice.
	Albums = append(Albums, a)
}

// Handles getAlbum
func GetAlbum(id string) (Album, error) {
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range Albums {
		if a.ID == id {
			return a, nil
		}
	}

	return Album{}, fmt.Errorf("album not found")
}
