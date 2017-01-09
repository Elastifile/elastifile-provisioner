package emanage

type DataStore struct {
	FreeSpace float64 `json:"free_space"`
	HostID    int     `json:"host_id"`
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
