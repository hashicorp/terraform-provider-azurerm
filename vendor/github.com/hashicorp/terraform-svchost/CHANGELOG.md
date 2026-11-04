## v0.2.0

#### Enhancements

Highly latent requests during service discovery are rare but happen many times per day. Service discovery now implements a hedged request strategy, where the service discovery request is repeated up to 7 times at 1500 ms intervals until the server responds.

## v0.1.1

The `disco.Disco` and `auth.CachingCredentialsSource` implementations are now safe for concurrent calls. Previously concurrent calls could potentially corrupt the internal cache maps or cause the Go runtime to panic.

## v0.1.0

#### Features:

- Adds hostname `Alias` method to service discovery, making it possible to interpret one hostname as another.

## v0.0.1

Initial release
