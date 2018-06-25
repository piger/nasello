# nasello

A very simple DNS proxy server capable of routing client queries to
different remote servers based on pattern matching.

This software relies heavily on [miekg](https://github.com/miekg)'s excellent
[DNS](https://github.com/miekg/dns) library.

Warning: this is alpha software and should be used with caution.

## Getting Started

### Getting nasello

The latest release is available at [Github][github-src]

[github-src]: https://github.com/piger/nasello

### Installing from source

You can install nasello with the standard `go get`:

	go install github.com/piger/nasello/cmd/nasello

This command will install a copy of `nasello` inside `$GOPATH/bin/`.

### Configuration format

The configuration file is a JSON document which must contains a
`filters` list with one or more `pattern` dictionaries; each `filter`
must contain a *FQDN* DNS name as the `pattern` and a list of one or
more remote DNS servers to forward the query to. For *reverse lookups*
the `in-addr.arpa` domain must be used in the pattern definition.

To use DNS-over-HTTPS see the example below.

The "." `pattern` specifies a default remote resolver.

### Example

The following configuration example specifies three `filters`:

- `*.example.com` will be resolved by OpenDNS (208.67.222.222, etc.)
- `10.0.24.*` will also be resolved by OpenDNS
- all the other queries will be forwarded to Google DNS (8.8.8.8,
  etc.)

`nasello.json`:

    {
      "filters": [
        {
          "pattern": "example.com.",
          "addresses": [ "208.67.222.222", "208.67.220.220" ]
        },
        {
          "pattern": "25.1.10.in-addr.arpa.",
          "addresses": [ "10.1.25.1", "10.1.25.2" ]
        },
        {
          "pattern": "google.com.",
          "addresses": [ "8.8.8.8", "8.8.4.4" ]
        },
        {
          "pattern": ".",
          "addresses": [ "1.1.1.1:853" ],
          "protocol": "tcp-tls"
        }
      ]
    }


## License

nasello is under the MIT license. See the [LICENSE][license] file for
details.

[license]: https://github.com/piger/nasello/blob/master/LICENSE
