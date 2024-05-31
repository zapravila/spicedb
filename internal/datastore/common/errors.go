package common

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"

	log "github.com/zapravila/spicedb/internal/logging"
	core "github.com/zapravila/spicedb/pkg/proto/core/v1"
	"github.com/zapravila/spicedb/pkg/spiceerrors"
	"github.com/zapravila/spicedb/pkg/tuple"
)

// SerializationError is returned when there's been a serialization
// error while performing a datastore operation
type SerializationError struct {
	error
}

func (err SerializationError) GRPCStatus() *status.Status {
	return spiceerrors.WithCodeAndDetails(
		err,
		codes.Aborted,
		spiceerrors.ForReason(
			v1.ErrorReason_ERROR_REASON_SERIALIZATION_FAILURE,
			map[string]string{},
		),
	)
}

func (err SerializationError) Unwrap() error {
	return err.error
}

// NewSerializationError creates a new SerializationError
func NewSerializationError(err error) error {
	return SerializationError{err}
}

// CreateRelationshipExistsError is an error returned when attempting to CREATE an already-existing
// relationship.
type CreateRelationshipExistsError struct {
	error

	// Relationship is the relationship that caused the error. May be nil, depending on the datastore.
	Relationship *core.RelationTuple
}

// GRPCStatus implements retrieving the gRPC status for the error.
func (err CreateRelationshipExistsError) GRPCStatus() *status.Status {
	if err.Relationship == nil {
		return spiceerrors.WithCodeAndDetails(
			err,
			codes.AlreadyExists,
			spiceerrors.ForReason(
				v1.ErrorReason_ERROR_REASON_ATTEMPT_TO_RECREATE_RELATIONSHIP,
				map[string]string{},
			),
		)
	}

	relationship := tuple.ToRelationship(err.Relationship)
	return spiceerrors.WithCodeAndDetails(
		err,
		codes.AlreadyExists,
		spiceerrors.ForReason(
			v1.ErrorReason_ERROR_REASON_ATTEMPT_TO_RECREATE_RELATIONSHIP,
			map[string]string{
				"relationship":       tuple.StringRelationshipWithoutCaveat(relationship),
				"resource_type":      relationship.Resource.ObjectType,
				"resource_object_id": relationship.Resource.ObjectId,
				"resource_relation":  relationship.Relation,
				"subject_type":       relationship.Subject.Object.ObjectType,
				"subject_object_id":  relationship.Subject.Object.ObjectId,
				"subject_relation":   relationship.Subject.OptionalRelation,
			},
		),
	)
}

// NewCreateRelationshipExistsError creates a new CreateRelationshipExistsError.
func NewCreateRelationshipExistsError(relationship *core.RelationTuple) error {
	msg := "could not CREATE one or more relationships, as they already existed. If this is persistent, please switch to TOUCH operations or specify a precondition"
	if relationship != nil {
		msg = fmt.Sprintf("could not CREATE relationship `%s`, as it already existed. If this is persistent, please switch to TOUCH operations or specify a precondition", tuple.StringWithoutCaveat(relationship))
	}

	return CreateRelationshipExistsError{
		fmt.Errorf(msg),
		relationship,
	}
}

var (
	portMatchRegex  = regexp.MustCompile("invalid port \\\"(.+)\\\" after host")
	parseMatchRegex = regexp.MustCompile("parse \\\"(.+)\\\":")
)

// RedactAndLogSensitiveConnString elides the given error, logging it only at trace
// level (after being redacted).
func RedactAndLogSensitiveConnString(ctx context.Context, baseErr string, err error, pgURL string) error {
	if err == nil {
		return errors.New(baseErr)
	}

	// See: https://github.com/jackc/pgx/issues/1271
	filtered := err.Error()
	filtered = strings.ReplaceAll(filtered, pgURL, "(redacted)")
	filtered = portMatchRegex.ReplaceAllString(filtered, "(redacted)")
	filtered = parseMatchRegex.ReplaceAllString(filtered, "(redacted)")
	log.Ctx(ctx).Trace().Msg(baseErr + ": " + filtered)
	return fmt.Errorf("%s. To view details of this error (that may contain sensitive information), please run with --log-level=trace", baseErr)
}
