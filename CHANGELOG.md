# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added 

## [0.0.7] - 2020-11-11

### Added 
- [Instalação Mosquito MQTT](https://www.instructables.com/Installing-MQTT-BrokerMosquitto-on-Raspberry-Pi/) mais um processo para automatizar (esse talvez possa deixar em docker)
- Criada Infrastructure MQTT seguindo [EMQ](https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang), [AWS](https://docs.aws.amazon.com/pt_br/whitepapers/latest/designing-mqtt-topics-aws-iot-core/mqtt-design-best-practices.html) e [Hive MQTT](https://www-hivemq-com.translate.goog/blog/mqtt-essentials-part-4-mqtt-publish-subscribe-unsubscribe/?_x_tr_sl=auto&_x_tr_tl=pt&_x_tr_hl=pt-BR&_x_tr_pto=nui)
- Domain robotAggregator refactor para Responsabilidade unica de metodos
- Domain para LCD e Sonar com logica de mensageria
- EnvVars e mais alguns outros fixes

## [0.0.6] - 2020-11-02

### Added

- [Instalação Gocv](https://gocv.io/getting-started/linux/)
- Adicionando uso simples de [Gocv](https://github.com/willmendil/golang_tutorials/blob/master/tutorial_1/main.go)/OpenCv para video stream
- Display LCD na linha 1 exibe ip e porta de rede da camera caso online)
- Env vars de uso da camera
- Acertos no keyboard e Wrapper de keyevent (ainda não está 100%)

## [0.0.5] - 2020-10-30

### Added

- Perifericos (input/output) podem ser habilitados por envVars
    - Montagem minima apenas com motores (ou qualquer outra combinação) agora pode ser configurada nas envs
    - Sem a necessidade de montar o projeto completo. Habilite apenas os recursos que vai plugar no esquema
- Adotando arquitetura hexagonal, desacoplando lib goBot com wrappers (problemas de desacoplamento com o keyboard)
- Série de pequenas correções por conta de conflitos (ainda existem dois repos, preciso remover do 'MatrixReality')

## [0.0.4] - 2020-10-24

### Added

- Usando viper para envVars
- Segregando envVars por modulos
- Corrigindo entradas do README

## [0.0.3] - 2020-10-23

### Added

- Comportamento de parada ao detectar obstáculo no sonar frontal do centro
- Inicio da documentação (CHANGELOG e melhorias no README.md, com mais commits do que gostaria de ter feito :| )
- Adicionado aos docs o esquema fritzing e arquivos de apoio a construção
- Adicionada pasta de scripts com o arquivo .ino do SonarSet


## [0.0.2] - 2020-10-23

### Added

- Mudança de repositório. Saindo do repo de estudos
- Refactor para arquitetura de periféricos, usando wrapper para gobot na maioria de inputs/outputs
- ['Driver' para Arduino](https://github.com/hybridgroup/gobot/blob/a8f33b2fc012951104857c485e85b35bf5c4cb9d/drivers/i2c/README.md)

[0.0.7]: https://github.com/jtonynet/autogo/compare/v0.0.6...v0.0.7
[0.0.6]: https://github.com/jtonynet/autogo/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/jtonynet/autogo/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/autogo/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/autogo/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/autogo/releases/tag/v0.0.2
