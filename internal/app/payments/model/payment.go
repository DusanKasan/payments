package model

import (
	errors "golang.org/x/xerrors"
	"regexp"
)

const ErrIdDoesNotExist = err("id does not exist")
const ErrIdAlreadyExist = err("id already exist")

type err string

func (e err) Error() string {
	return string(e)
}

type Payment struct {
	Type           string                 `json:"type"`
	ID             string                 `json:"id"`
	Version        int16                  `json:"version"`
	OrganisationID string                 `json:"organisation_id"`
	Attributes     map[string]interface{} `json:"attributes"`
}

func (p Payment) Validate() error {
	if p.Version != 0 {
		return errors.New("version must be 0")
	}

	if p.Type != "Payment" {
		return errors.New("type must be 'Payment'")
	}

	if !isValidUUID(p.ID) {
		return errors.New("id is not an uuid")
	}

	if !isValidUUID(p.OrganisationID) {
		return errors.New("organisation_id is not an uuid")
	}

	return nil
}

// TODO: could be sped up by using something like https://github.com/google/uuid/blob/dec09d789f3dba190787f8b4454c7d3c936fed9e/uuid.go#L41
func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}