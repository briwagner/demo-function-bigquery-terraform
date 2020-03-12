package bigquery_etl

type Person struct {
	Name    string `json:"name" csv:"Name" bigquery:"name"`
	Email   string `json:"email" csv:"Mail" bigquery:"email"`
	Contact string `json:"contact" csv:"Contact" bigquery:"contact"`
	Pet     string `json:"pet" csv:"Pet" bigquery:"-"`
}
