States
Actions/transitions

Start
- display form(source ns, dst address, dst port)

DestinationIsIP
- display massage
- go back/continue buttons

DestinationIsFQDN
- check network policies in source ns

NetworkPolicyExists
    if DestinationIsFQDN
    - perform DNS resolution
    - if error   - display error message
    - if success - go to IP check
    if DestinationIsIP
    - skip dns resolution, go to IP check

NetworkPolicyDoesNotExist
- display error message + wiki links

IPcheck
if isIPpublicNonIngress
- display relevant message
if isIPcoreIngress
- display message
if isIPprivateNonIngress
- display message
