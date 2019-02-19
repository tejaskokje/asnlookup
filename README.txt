asnlookup - Utility to lookup ASN from pre-defined list for a given IPv4 or IPv6 address

Author: Tejas Kokje (tejas.kokje@gmail.com)

Description
-----------

asnlookup is an utility to lookup ASN for IPv4 or IPv6 address. It expects only one IPv4 or IPv6 address as an target.

By default it builds it database by fetching information from http://lg01.infra.ring.nlnog.net/table.txt. If CONFIG_FILE_PATH environment variable is defined, utility uses value of CONFIG_FILE_PATH environment variable as file name and builds its database. 

This utility will report Subnets, CIDR  & ASN sorted by CIDR prefix.

Usage
-----

asnlookup <IPv4 or IPv6 address>

Example:
    
    asnlookup 8.8.8.8
    
    asnlookup 2001:db8:0:b::2a:1a

Note that for IPv6, result IPv6 CIDR block will always be displayed in uncompressed format.

Design
------

Internally, this utility uses binary trie to store subnet, CIDR & ASN information for efficient lookup. Trie is built when addresses are parsed from either default URL or configuration file. For each address (IPv4 or IPv6), it starts looking at higher order bits. If the bit is 0, it goes left. If bit is 1, it goes right. When CIDR prefix bits end, the value (subnet, CIDR prefix length, asn) is stored at that node.

While doing lookup, it just walks through the trie using bits of target IP address. If any trie node along the way has any (subnet, CIDR, asn) values set, it stores them in a list. Finally, before printing, it sorts the list using CIDR prefix length. This is very similar to how routers perform route lookup except that this utility returns all matched entries instead of just longest prefix.

It can be observed that using binary trie to store information might not be optimal, especially if number of subnets are sparse. Path compression using Patricia trie or multibit trie might be more efficient. However, this utility only takes one IP address as target for lookup. Hence, cost to optimize for path compression or finding optimal number of bits for multibit trie is be more than just searching for target IP address. It could be beneficial to use efficient trie structure if this utility is enhanced to search for multiple IP addresses or is turned into a service accepting requests for ASN lookups.

Even though asnlookup can do lookup for IPv4 or IPv6 address, it only maintains only one trie based on target IP address type.

Implementation
--------------

This utility is implemented using Go programming language. It only uses packages available in Go standard library and does not import any external packages. For IPv4 and IPv6 address parsing & processing, it does not rely on Go's net package. All IPv4 & IPv6 parsing functions are implemented inside this utility.

Common interface (called IPAaddress) is implemented by IPv4 & IPv6 address types. This allows for binary trie to be address type agnostic. Logic to store and retrieve information from trie remains same for both type of addresses. Only difference is that trie size for IPv4 is smaller than for IPv6.

Go unit tests are implemented for all files. 

Design Decisions
---------------
Following were the focus of design for this utility

1) Efficiency - Binary trie used to store and retrieve information is time and space efficient. However, it does have penalty of doing lot of memory lookups (32 maximum for IPv4 & 128 maximum for IPv6). Given the scope of the utility to process only one target IP address, binary trie seems to be optimal.

2) Readability - Readability suffers a bit in IPv4 and IPv6 address parsing & processing code. There is lot of bit shifting & string processing which might not be obvious just by reading code. Overall, other parts of the code including trie logic & configuration parsing are very readable.

3) Mantainability - The code is spread out into relevant files with interface & functions used at appropriate places. For e.g. If Go's IP address primitive from net package is allowed to be used, existing IPv4 & IPv6 parsing function body could just be call to relevant API (with some processing) from net package.

4) Debuggability - Code to store and retrive information from trie is fairly small. There is also utility function to dump trie nodes. This makes code to fairly debuggable.

Build Instructions
------------------

To build this utility, simply invoke following command from the package

#make build

This should build asnlookup binary in current working directory.

If you are using "go" tool natively to build, you will have to ensure that current directory is in GOPATH. 


