#### Filtering go http transport and proxy handler

See [client](filterclient/main.go) and [proxy](filterproxy/main.go) examples.

#### Known issues

- Uses socket and connect syscall. Not test and probably does not work on Windows etc.
- No timeout config
- Does not do IPv6 happy eyeballs

#### License

filtertransport is licensed under the MIT license. See [LICENSE](LICENSE) for the full license text.
