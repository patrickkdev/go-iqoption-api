# Go API Wrapper for Quadcore Brokers (Like IQ Option)

This is an unofficial Go implementation of the existing Python API used for Quadcore-based brokers such as IQ Option, Exnova, and others. It was originally a private project and is being open-sourced on December 2 2025.

The wrapper worked reliably when it was actively maintained, but the upstream API changes often and requires ongoing attention. Since this is based on reverse engineering rather than an official API, breakage is expected.

The implementation was rewritten from the Python version with the goal of matching its structure and behavior as closely as possible. Contributions are welcome if you want to help keep it updated.

This library powered a full copy-trading platform with rooms, invitations, live updates, automatic trading mode, signal-list mode, and other tooling. Development slowed down only because of the time needed to keep up with API changes.

Feel free to use this as a reference, a starting point, or a building block for your own projects.

## Status

The current state is uncertain, as it has been a while since this was last used in production. It should still work for the most part, and there likely aren’t too many changes required. Expect occasional breakage when Quadcore updates their protocols.

## Supported Brokers

* IQ Option
* Exnova
* Any other broker powered by Quadcore

## Contributing

Pull requests and issues are welcome. If you have protocol updates, maintenance fixes, or improvements, feel free to send them.

## Disclaimer

This api is intended to be an open source project to communicate with Quadcore websocket. This is not a official repository, it means it is maintained by community.

Due to the large amount of scammers that have appeared in the market, it is recommended that you DO NOT enter your password into an unknown exe or robot site that operates on iqoption because many of those have stolen people's passwords so be careful. It's best if you develop your robot or hire someone you trust. [Look here](https://patrick.makztech.com)

Esta API é destinada a ser um projeto de código aberto para se comunicar com o websocket da Quadcore. Este é um repositório não oficial, significa que é mantido pela comunidade.

Devido a grande quantidade de golpistas que tem aparecido no mercado, recomenda-se que você NÃO inserir sua senha em exe ou sites de robo desconhecidos que opera na iqoption porque muitos desses tem roubado as senhas das pessoas então tomem cuidado. O melhor é você desenvolver seu robô ou contratar alguém de confiança. [Encontre aqui](https://patrick.makztech.com)
