# Changelog

## 1.0.0 (2024-01-05)


### Features

* Add Close function on top client level ([f319b46](https://github.com/d-strobel/gowindows/commit/f319b4605eaa936922761e036fc6062408fbf7cc))
* add context to ssh execution ([645caab](https://github.com/d-strobel/gowindows/commit/645caabc4fc728ebfede9df16e7ed6739d63a9b0))
* add custom error type for parser ([4798e03](https://github.com/d-strobel/gowindows/commit/4798e034be709c90f562e64c93b971ab13a6ce24))
* add custom errors to windows functions ([5620368](https://github.com/d-strobel/gowindows/commit/56203686d9b6b8035ed002d55c3b675db9e6c3ab))
* add errorhandling to run function ([f7f443e](https://github.com/d-strobel/gowindows/commit/f7f443eda46fef926cc40d49f1c696e80779ef04))
* add json unmarshal for WinTime ([15e5fad](https://github.com/d-strobel/gowindows/commit/15e5fade06a3a37088c6885d9f1b0a31261025ff))
* add kerberos auth to winrm ([e4fde80](https://github.com/d-strobel/gowindows/commit/e4fde80dda2defa52fac5cfd0af8950e252d21d9))
* add local group member functions ([5d352cc](https://github.com/d-strobel/gowindows/commit/5d352cca7acba2b178f5e33bf69e67abe8794f3a))
* add local user create ([e7f2ef0](https://github.com/d-strobel/gowindows/commit/e7f2ef0979b4694b8935359ac0c89781cdfe3f16))
* add local user delete ([a955dc5](https://github.com/d-strobel/gowindows/commit/a955dc50cae49c433130cf2ee3e30b627b7e6bd1))
* add local user list function ([93f3ec0](https://github.com/d-strobel/gowindows/commit/93f3ec0ca18d6d93f4bb468aefd688e07810cc0d))
* add local user read ([a46c12a](https://github.com/d-strobel/gowindows/commit/a46c12abbe0b6c6524ae9f3f1d100abd7f8e89e2))
* add local user update function ([28524f6](https://github.com/d-strobel/gowindows/commit/28524f6ac83f70d0318fb1b3b2e8c8b24e785913))
* add local-GroupCreate function ([99ff835](https://github.com/d-strobel/gowindows/commit/99ff835af117d54bbc48b404c5cfd3cdacd951e2))
* add local-GroupDelete function ([7689f2f](https://github.com/d-strobel/gowindows/commit/7689f2facf4e51f282bf0ca54ca558fb31e6a929))
* add local-GroupList function ([8ae09ff](https://github.com/d-strobel/gowindows/commit/8ae09ff60b2df523d6dc2c218919d4e3b6146fb6))
* add local-GroupUpdate function ([864a46b](https://github.com/d-strobel/gowindows/commit/864a46b9029c663b11221275480ec72beabdd8c5))
* add parser package for inputs and outputs ([526bf85](https://github.com/d-strobel/gowindows/commit/526bf85eb763ab3531db59f36e90c06c4ec9eee4))
* add private key authentication method ([145586c](https://github.com/d-strobel/gowindows/commit/145586c7457615345e6267f0235f9b56c1cc21f0))
* add ssh host key callback check ([f085555](https://github.com/d-strobel/gowindows/commit/f0855553a155e5e0ee5ca0f9789147b3d39aa03e))
* Add winrm client ([445321a](https://github.com/d-strobel/gowindows/commit/445321ab4f6ff2126ebeaee81dd52b90dd989b31))
* change function return values ([8f2b39f](https://github.com/d-strobel/gowindows/commit/8f2b39f68d92c60ee7e9cc666825c8f8cda27959))
* Change package name and refactor ([04a6f86](https://github.com/d-strobel/gowindows/commit/04a6f86462ce350fd5ed66d15294f3995504fa70))
* compress json output ([4261f6b](https://github.com/d-strobel/gowindows/commit/4261f6bd6ab2cc8d128e648bd369810205407ed6))
* compress json output ([8515df8](https://github.com/d-strobel/gowindows/commit/8515df8cc34b8d5747eff82602d84997a72f755c))
* implement custom errors ([ac1ccf1](https://github.com/d-strobel/gowindows/commit/ac1ccf1d870826e3078e5fb0e2cd21e6fd46d6aa))
* replace error with custom error ([a26778e](https://github.com/d-strobel/gowindows/commit/a26778eddd3091972aee39656ba36e6c9774635f))
* replace errors with custom errors ([ecd0494](https://github.com/d-strobel/gowindows/commit/ecd049429c326e90eee80efdf389906e6001cb6c))
* replace errors with custom errors ([3466dac](https://github.com/d-strobel/gowindows/commit/3466dac0158e5057ab27bbdf5cf295864d0a4a2f))
* replace errors with custom errors ([c59fa92](https://github.com/d-strobel/gowindows/commit/c59fa92148cf0716f3967e9cb3f49628c9c71bf3))
* winrm default variables depends on tls ([d54347c](https://github.com/d-strobel/gowindows/commit/d54347c812bdc8ea7e7e13b14af76668d4e073a8))


### Bug Fixes

* cant read group when name has spaces ([d7be6bc](https://github.com/d-strobel/gowindows/commit/d7be6bc2ba9d13c63de054172d843765f90a8318))
* clixml normalizing ([97a01c0](https://github.com/d-strobel/gowindows/commit/97a01c06eceedbfe7fe45342f2337b26265e40d0))
* description param cant be removed with update ([4ae5be7](https://github.com/d-strobel/gowindows/commit/4ae5be70166bcb01990e8f29bfbc6035d616d480))
* group name to wildcard leads to error ([9656e61](https://github.com/d-strobel/gowindows/commit/9656e6102a48658da6b3f2798e916c95b7bb868d))
* local.UserUpdate function ([38d7580](https://github.com/d-strobel/gowindows/commit/38d7580d57de48d5844e82e3e1dfb5151a7b817c))
* ssh error handling ([c40c785](https://github.com/d-strobel/gowindows/commit/c40c785e2ff42dda7be84d195aab3d729968ffc2))
* ssh returns error if stdout and stderr are empty ([5785490](https://github.com/d-strobel/gowindows/commit/57854903dd6b93c11e2dade6908fe428f0b070a4))
* ssh stdout and stderr are empty ([c5452c0](https://github.com/d-strobel/gowindows/commit/c5452c08b9fc1bf2e5ad8b38f52f38779f8ed8f8))
* stderr handling ([6a4aafb](https://github.com/d-strobel/gowindows/commit/6a4aafb772e4ffcf751f271d6681928ae5b973b2))


### Reverts

* go build version 1.19 and 1.20 ([e5b9187](https://github.com/d-strobel/gowindows/commit/e5b9187508244cf5299cf1e4647cc9cb97d4ccb3))
* move commitlint into build action ([917557e](https://github.com/d-strobel/gowindows/commit/917557e8d7595ca192c6e33ae3b1d44bffd37f3d))
