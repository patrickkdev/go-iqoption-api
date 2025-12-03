# Go API Wrapper for Quadcore Brokers (Like IQ Option)

This is an unofficial Go implementation of the existing Python API used for Quadcore-based brokers such as IQ Option, Exnova, and others. It was originally a private project and is being open-sourced on December 2.

The wrapper used to work flawlessly, but the upstream API changes often and requires ongoing maintenance. Since this uses reverse engineering and not an oficial API, breakage is expected from time to time.

The implementation was rewritten from the Python version with the goal of keeping the workflow and behavior as close as possible. If you want to help keep it updated and healthy, contributions are welcome.

This library powered a complete copy-trading platform: rooms, invite system, live updates, automatic trading mode, signal-list mode, and other tooling. The only reason it stalled is lack of time to keep chasing API changes.

Use it as a reference, a starting point, or a building block for your own projects.

## Status

Not sure, it has been a while since I last used this, it should still work mostly. Probably not many things to update. Expect occasional breakage when Quadcore changes their protocols.

## Supported Brokers

* IQ Option
* Exnova
* Any other broker powered by quadcore

## Contributing

Pull requests and issues are welcome. If you have protocol updates or bug fixes, feel free to send them.
