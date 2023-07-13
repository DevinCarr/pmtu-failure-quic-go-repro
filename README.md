This repo tests a failure case that happens when the PMTU is reduced during an active connection in quic-go.

The quic-go library does use an MTU discovery process after the handshake is completed between the client and the
server. You can see this code here: https://github.com/quic-go/quic-go/blob/master/mtu_discoverer.go. This code does
not handle the case that the PMTU is lowered below the "current" byte size. So if you have a long-lived connection
between a server and a client that is sending ping frames but no other payloads, a PMTU reduction will not cause any
issues until either side attempts to send a payload that would normally fill the QUIC packet size. Afterwhich, a payload
that would exceed the currently calculated PMTU will be dropped by along the path where the reduction occurred.

Normal operation to trigger the failure can be done as followed:
 0. Set the `addr` to two IPs that can be reached between two separate machines for both `client.go` and `server/server.go`.

    This needs to be done between two machines as I believe there is some WriteTo/ReadFrom's that will not
    reproduce the problem; there needs to be a PMTU that can be adjusted.
    A network interface shared between both the client and the server will use a 65K MTU. Similarly, this will not work 
    if you attempt to create a dummy link as the kernel will again just pipe between the two processes instead of assigning 
    to the network interface the adjusted MTU.
    Addtionally, this cannot be done on localhost as the kernel treats it differently than a network interface and 
    typically has a 65536 MTU that should (loose guideline) not be adjusted; in either case, localhost doesn't exhibit
    this problem.

 1. Start echo server (`go run server/server.go`)

 2. Start the client (`go run client.go`)

 3. Hit enter on the client process to send the first payload normally

 4. Adjust MTU on interface link that is used: `sudo ip link dev <dev> set mtu 1350`

    If between two machines, either side should work, but I did my testing by adjusting the client's network link MTU.
    1439 is the typical value reached on my machine during the PMTU discovery process, your machine may vary. The
    goal is to adjust the MTU to be at least one byte less than the current PMTU set by quic-go; 1350 worked reliably for
    my testing.

 5. Hit enter on the client process again to see the payload fail to send and timeout.