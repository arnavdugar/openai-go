// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/arnavdugar/openai-go/internal/apiform"
	"github.com/arnavdugar/openai-go/internal/apijson"
	"github.com/arnavdugar/openai-go/internal/requestconfig"
	"github.com/arnavdugar/openai-go/option"
	"github.com/arnavdugar/openai-go/packages/respjson"
	"github.com/arnavdugar/openai-go/shared/constant"
)

// UploadPartService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewUploadPartService] method instead.
type UploadPartService struct {
	Options []option.RequestOption
}

// NewUploadPartService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewUploadPartService(opts ...option.RequestOption) (r UploadPartService) {
	r = UploadPartService{}
	r.Options = opts
	return
}

// Adds a
// [Part](https://platform.openai.com/docs/api-reference/uploads/part-object) to an
// [Upload](https://platform.openai.com/docs/api-reference/uploads/object) object.
// A Part represents a chunk of bytes from the file you are trying to upload.
//
// Each Part can be at most 64 MB, and you can add Parts until you hit the Upload
// maximum of 8 GB.
//
// It is possible to add multiple Parts in parallel. You can decide the intended
// order of the Parts when you
// [complete the Upload](https://platform.openai.com/docs/api-reference/uploads/complete).
func (r *UploadPartService) New(ctx context.Context, uploadID string, body UploadPartNewParams, opts ...option.RequestOption) (res *UploadPart, err error) {
	opts = append(r.Options[:], opts...)
	if uploadID == "" {
		err = errors.New("missing required upload_id parameter")
		return
	}
	path := fmt.Sprintf("uploads/%s/parts", uploadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// The upload Part represents a chunk of bytes we can add to an Upload object.
type UploadPart struct {
	// The upload Part unique identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the Part was created.
	CreatedAt int64 `json:"created_at,required"`
	// The object type, which is always `upload.part`.
	Object constant.UploadPart `json:"object,required"`
	// The ID of the Upload object that this Part was added to.
	UploadID string `json:"upload_id,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		UploadID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UploadPart) RawJSON() string { return r.JSON.raw }
func (r *UploadPart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UploadPartNewParams struct {
	// The chunk of bytes for this Part.
	Data io.Reader `json:"data,omitzero,required" format:"binary"`
	paramObj
}

func (r UploadPartNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}
