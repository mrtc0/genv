package mocks

const (
	StsGetRoleCredentialsAccountId       = "111122223333"
	StsGetRoleCredentialsRoleName        = "SSORole"
	StsGetRoleCredentialsAccessKeyId     = "ASIAI44QH8DHBEXAMPLE"
	StsGetRoleCredentialsSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	StsGetRoleCredentialsSessionToken    = "IQoJb3JpZ2luX2IQoJb3JpZ2luX2IQoJb3JpZ2luX2IQoJb3JpZ2luX2IQoJb3JpZVERYLONGSTRINGEXAMPLE"

	MockStsGetCallerIdentityValidResponseBody = `
<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>arn:aws:iam::222222222222:user/Alice</Arn>
    <UserId>AKIAI44QH8DHBEXAMPLE</UserId>
    <Account>222222222222</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`
	MockStsGetCallerIdentityValidAssumedRoleResponseBody = `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>arn:aws:sts::555555555555:assumed-role/role/AssumeRoleSessionName</Arn>
    <UserId>ARO123EXAMPLE123:AssumeRoleSessionName</UserId>
    <Account>555555555555</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>
`

	MockStsGetRoleCredentialsValidResponseBodyTemplate = `
{
   "roleCredentials": { 
      "accessKeyId": "%s",
      "expiration": %d,
      "secretAccessKey": "%s",
      "sessionToken": "%s"
   }
}
`
)
