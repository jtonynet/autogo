# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added 
- Wrapper de keyevent
- adicionando uso simples de Gocv para video stream

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

[0.0.5]: https://github.com/jtonynet/autogo/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/autogo/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/autogo/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/autogo/releases/tag/v0.0.2
