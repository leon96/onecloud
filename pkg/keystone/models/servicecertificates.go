package models

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/apis"
	api "yunion.io/x/onecloud/pkg/apis/identity"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/validators"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type SServiceCertificateManager struct {
	db.SStandaloneResourceBaseManager
}

var ServiceCertificateManager *SServiceCertificateManager

func init() {
	ServiceCertificateManager = &SServiceCertificateManager{
		SStandaloneResourceBaseManager: db.NewStandaloneResourceBaseManager(
			SServiceCertificate{},
			"servicecertificates_tbl",
			"servicecertificate",
			"servicecertificates",
		),
	}
	ServiceCertificateManager.SetVirtualObject(ServiceCertificateManager)
}

type SServiceCertificate struct {
	db.SStandaloneResourceBase
	db.SCertificateResourceBase

	CaCertificate string `create:"optional" list:"admin"`
	CaPrivateKey  string `create:"optional" list:"admin"`
}

func (man *SServiceCertificateManager) ValidateCreateData(ctx context.Context, userCred mcclient.TokenCredential, ownerId mcclient.IIdentityProvider, query jsonutils.JSONObject, data *jsonutils.JSONDict) (*jsonutils.JSONDict, error) {
	v := validators.NewCertKeyValidator("certificate", "private_key")
	if err := v.Validate(data); err != nil {
		return nil, err
	}
	data = v.UpdateCertKeyInfo(ctx, data)

	// validate ca cert key
	v = validators.NewCertKeyValidator("ca_certificate", "ca_private_key")
	if err := v.Validate(data); err != nil {
		return nil, err
	}

	input := apis.StandaloneResourceCreateInput{}
	err := data.Unmarshal(&input)
	if err != nil {
		return nil, httperrors.NewInternalServerError("unmarshal StandaloneResourceCreateInput fail %s", err)
	}
	input, err = man.SStandaloneResourceBaseManager.ValidateCreateData(ctx, userCred, ownerId, query, input)
	if err != nil {
		return nil, err
	}
	data.Update(jsonutils.Marshal(input))
	return data, nil
}

func (cert *SServiceCertificate) ValidateUpdateData(
	ctx context.Context, userCred mcclient.TokenCredential,
	query jsonutils.JSONObject, data *jsonutils.JSONDict,
) (*jsonutils.JSONDict, error) {
	v := validators.NewCertKeyValidator("certificate", "private_key")
	if err := v.Validate(data); err != nil {
		return nil, err
	}
	data = v.UpdateCertKeyInfo(ctx, data)

	// validate ca cert key
	v = validators.NewCertKeyValidator("ca_certificate", "ca_private_key")
	if err := v.Validate(data); err != nil {
		return nil, err
	}

	updateData := jsonutils.NewDict()
	if name, err := data.GetString("name"); err == nil {
		updateData.Set("name", jsonutils.NewString(name))
	}

	if desc, err := data.GetString("description"); err == nil {
		updateData.Set("description", jsonutils.NewString(desc))
	}

	input := apis.StandaloneResourceBaseUpdateInput{}
	err := updateData.Unmarshal(&input)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	input, err = cert.SStandaloneResourceBase.ValidateUpdateData(ctx, userCred, query, input)
	if err != nil {
		return nil, errors.Wrap(err, "SVirtualResourceBase.ValidateUpdateData")
	}
	updateData.Update(jsonutils.Marshal(input))

	return updateData, nil
}

func (cert *SServiceCertificate) ToOutput() *api.CertificateDetails {
	return &api.CertificateDetails{
		SCertificateResourceBase: apis.SCertificateResourceBase{
			Certificate:             cert.Certificate,
			PrivateKey:              cert.PrivateKey,
			PublicKeyAlgorithm:      cert.PublicKeyAlgorithm,
			PublicKeyBitLen:         cert.PublicKeyBitLen,
			SignatureAlgorithm:      cert.SignatureAlgorithm,
			Fingerprint:             cert.Fingerprint,
			NotAfter:                cert.NotAfter,
			NotBefore:               cert.NotBefore,
			CommonName:              cert.CommonName,
			SubjectAlternativeNames: cert.SubjectAlternativeNames,
		},
		CertName:      cert.Name,
		CertId:        cert.Id,
		CaCertificate: cert.CaCertificate,
		CaPrivateKey:  cert.CaPrivateKey,
	}

}
