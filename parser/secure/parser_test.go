package secret

import (
	"fmt"
	"testing"

	"github.com/franela/goblin"
)

func Test_Secure(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Secure yaml", func() {

		priv, _ := decodePrivateKey(fakePriv)
		pub := priv.PublicKey
		pem := encodePrivateKey(priv)

		fmt.Println(encrypt(mapYaml, &pub))

		g.It("Should encrypt and decrypt", func() {
			plain := checksumYaml
			encrypted, err := encrypt(plain, &pub)
			g.Assert(err == nil).IsTrue()
			decrypted, err := decrypt(encrypted, priv)
			g.Assert(err == nil).IsTrue()
			g.Assert(plain).Equal(string(decrypted))
		})

		g.It("Should decrypt a yaml and verify secrets", func() {
			secure, err := Parse(dockerEnc, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(len(secure.Runtime)).Equal(2)

			var index = secure.Registry[0]
			var docker = secure.Runtime[0]
			var amazon = secure.Runtime[1]
			if secure.Runtime[0].Image[0] == "docker" {
				docker = secure.Runtime[0]
				amazon = secure.Runtime[1]
			} else {
				docker = secure.Runtime[1]
				amazon = secure.Runtime[0]
			}

			g.Assert(secure.Checksum).Equal("fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b")
			g.Assert(docker.Image).Equal([]string{"docker"})
			g.Assert(docker.Event).Equal([]string{"push", "pull_request"})
			g.Assert(docker.Data["username"]).Equal("octocat")
			g.Assert(docker.Data["password"]).Equal("pa55word")

			g.Assert(amazon.Image).Equal([]string{"s3", "aws"})
			g.Assert(amazon.Data["access_key"]).Equal("AKID1234567890")
			g.Assert(amazon.Data["secret_access_key"]).Equal("MY-SECRET-KEY")

			g.Assert(index.Hostname).Equal("index.docker.io")
			g.Assert(index.Username).Equal("octocat")
			g.Assert(index.Password).Equal("pa55word")
			g.Assert(index.Email).Equal("octocat@github.com")
		})

		g.It("Should decrypt a yaml with empty secrets", func() {
			secure, err := Parse(checksumEnc, pem)
			g.Assert(err == nil).IsTrue()
			g.Assert(secure.Checksum).Equal("fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b")
			g.Assert(len(secure.Runtime)).Equal(0)
			g.Assert(len(secure.Registry)).Equal(0)
		})
	})
}

var mapYaml = `
checksum: fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b

runtime:
  docker:
    event: [ push, pull_request ]
    username: octocat
    password: pa55word
  amazon:
    image: [ s3, aws ]
    access_key: AKID1234567890
    secret_access_key: MY-SECRET-KEY

registry:
  index.docker.io:
    username: octocat
    password: pa55word
    email: octocat@github.com
`

/*
registry:
  username: foo
  password: bar
  email: baz@foo.com
*/

var checksumYaml = `
checksum: fa4d4048a6bd1a94f2775039ecf29b812d9cfe6b
`

var dockerEnc = `eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.sIOM0ib607gH0t5_uQoUk2Quo98vDvEupMLDw4KVqYB_NJad9WvbID99V7p5KDg6ULGGLRFISld_INN85dI_Ocxl6eDyqRXPmGhm4UJ3utuUJ6xrgFYqooSRvnZpAYl30NRVzvlOTKdalR0nia_owAqjjOxKG5GU3BKiyv6jHmTm-wiy7pPpuPdDZ2TFyrxpBHgaWzPjyqrBAaT7DgTLCh-kTZPF5cSLSCRMT0AuOUcKrg82kmztIvmdEI2XL_e_IBKZLV6R-Eu5wNeLJ7OSm35VpG7zH4882oMLiImpmmGFdU4pYp6jk7lNsf9Xwo4oD4xJaKwEV6zBuiC9IXqtLg.z8l3DeLADxBvXN0_.hrxREzBUJfQPDKJgVaJhB1XWejDqF2MH4ztHeJcda83qEJeOcdgpFTJiwu7osV_p_J43-GfMzm9iWCmEQIQM0U9JZ1wiukG7lYY1yMMuzh_zjZvzwGPcnUJ4gK64JtvTlMoJHlhpMIRIrpmGDVKHZzgplKGQBBYUsei8YrNSGnyTCPAbDekMUmmJrZizHnlA6WOKqy_Ne8zAuLJHGEbgq-dpeMUsSozJE7Xus3CXRHXMWAOhdyE_jpgEvzXUwFNx8GsHK-CWVdXhoZOcjt70hTlP-puntWJkF1DDtAd1VeT9ooIH4dG4T5APKk_PKTEzWvlUUMABseQlGucrCv9QeympJIcsMg_oMeYZBj1dmRioIlegLJ5Z0xzq4jCpOMQU7f-Tu60gI8mWFa_5bkHVsRZTVhnAp5nbvzMisJ-k-09tpE97zSdTavbp3wV0GnXqBEfj-LuVGcE2F13ecfPejpCL-kge.Y94i9PVttepV0aBdiygdBA`

var checksumEnc = `eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.lduGCINc5DVUD6hi3UHzWNkuKLlsrXudLfrktD3gOI36J6r58DlMAGcUNLfgSAU4v0kf9L407EkB4dtqwXbkhNeihNw69BbYa94QQ1H2uW3BNQbPq-1JeeZLAU6dXmkRXZur4KGNuWls4tMqd-Z9OyRSCBzogDzMf2JGJ-eLSL63zhBCzKGwQ6yE1N6cZsS2NN0-1BEgrAk-dC76motQvcRTHmiosADrEGUAM6xy-LQSkcC8DImpXajv-AFlFv5F4BtFBg9e7MLrVishwZAFKq-lexWLRqlcf7xqgU5GVt6_3VtuoWtVIyUFP4ZnM0KFScKrG6zsd1h7G5_zSf9AMA.YaL_NLt5Ei_7BEeo.7igBRj8A-EfvsT4VafSBCi_68_lelDwcANbtePmZENuuxEaLSRyfMsawKY0Oyc9DdYCsNA.Np3a7xQRMNHK1z6Nb8oGXg`

var fakePriv = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA71FaA+otDak2rXF/4h69Tz+OxS6NOWaOc/n7dinHXnlo3Toy
ZzvwweJGQKIOfPNBMncz+8h6oLOByFvb95Z1UEM0d+KCFCCutOeN9NNMw4fkUtSZ
7sm6T35wQUkDOiO1YAGy27hQfT7iryhPwA8KmgZmt7toNNf+WymPR8DMwAAYeqHA
5DIEWWsg+RLohOJ0itIk9q6Us9WYhng0sZ9+U+C87FospjKRMyAinSvKx0Uan4ap
YGbLjDQHimWtimfT4XWCGTO1cWno378Vm/newUN6WVaeZ2CSHcWgD2fWcjFixX2A
SvcvfuCo7yZPUPWeiYKrc5d1CC3ncocu43LhSQIDAQABAoIBAQDIbYKM+sfmxAwF
8KOg1gvIXjuNCrK+GxU9LmSajtzpU5cuiHoEGaBGUOJzaQXnQbcds9W2ji2dfxk3
my87SShRIyfDK9GzV7fZzIAIRhrpO1tOv713zj0aLJOJKcPpIlTZ5jJMcC4A5vTk
q0c3W6GOY8QNJohckXT2FnVoK6GPPiaZnavkwH33cJk0j1vMsbADdKF7Jdfq9FBF
Lx+Za7wo79MQIr68KEqsqMpmrawIf1T3TqOCNbkPCL2tu5EfoyGIItrH33SBOV/B
HbIfe4nJYZMWXhe3kZ/xCFqiRx6/wlc5pGCwCicgHJJe/l8Y9OticDCCyJDQtD8I
6927/j2NAoGBAPNRRY8r5ES5f8ftEktcLwh2zw08PNkcolTeqsEMbWAQspV/v+Ay
4niEXIN3ix2yTnMgrtxRGO7zdPnMaTN8E88FsSDKQ97lm7m3jo7lZtDMz16UxGmd
AOOuXwUtpngz7OrQ25NXhvFYLTgLoPsv3PbFbF1pwbhZqPTttTdg5so3AoGBAPvK
ta/n7DMZd/HptrkdkxxHaGN19ZjBVIqyeORhIDznEYjv9Z90JvzRxCmUriD4fyJC
/XSTytORa34UgmOk1XFtxWusXhnYqCTIHG/MKCy9D4ifzFzii9y/M+EnQIMb658l
+edLyrGFla+t5NS1XAqDYjfqpUFbMvU1kVoDJ/B/AoGBANBQe3o5PMSuAD19tdT5
Rnc7qMcPFJVZE44P2SdQaW/+u7aM2gyr5AMEZ2RS+7LgDpQ4nhyX/f3OSA75t/PR
PfBXUi/dm8AA2pNlGNM0ihMn1j6GpaY6OiG0DzwSulxdMHBVgjgijrCgKo66Pgfw
EYDgw4cyXR1k/ec8gJK6Dr1/AoGBANvmSY77Kdnm4E4yIxbAsX39DznuBzQFhGQt
Qk+SU6lc1H+Xshg0ROh/+qWl5/17iOzPPLPXb0getJZEKywDBTYu/D/xJa3E/fRB
oDQzRNLtuudDSCPG5wc/JXv53+mhNMKlU/+gvcEUPYpUgIkUavHzlI/pKbJOh86H
ng3Su8rZAn9w/zkoJu+n7sHta/Hp6zPTbvjZ1EijZp0+RygBgiv9UjDZ6D9EGcjR
ZiFwuc8I0g7+GRkgG2NbfqX5Cewb/nbJQpHPO31bqJrcLzU0KurYAwQVx6WGW0He
ERIlTeOMxVo6M0OpI+rH5bOLdLLEVhNtM/4HUFi1Qy6CCMbN2t3H
-----END RSA PRIVATE KEY-----
`
