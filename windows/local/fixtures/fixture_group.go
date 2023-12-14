package fixtures

const UsersGroup = `{
    "Description":  "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
    "Name":  "Users",
    "SID":  {
                "BinaryLength":  16,
                "AccountDomainSid":  null,
                "Value":  "S-1-5-32-545"
            },
    "PrincipalSource":  1,
    "ObjectClass":  "Group"
}`

const GroupList = `[
    {
        "Description":  "Administrators have complete and unrestricted access to the computer/domain",
        "Name":  "Administrators",
        "SID":  {
                    "BinaryLength":  16,
                    "AccountDomainSid":  null,
                    "Value":  "S-1-5-32-544"
                },
        "PrincipalSource":  1,
        "ObjectClass":  "Group"
    },
    {
        "Description":  "Users are prevented from making accidental or intentional system-wide changes and can run most applications",
        "Name":  "Users",
        "SID":  {
                    "BinaryLength":  16,
                    "AccountDomainSid":  null,
                    "Value":  "S-1-5-32-545"
                },
        "PrincipalSource":  1,
        "ObjectClass":  "Group"
    }
]`
