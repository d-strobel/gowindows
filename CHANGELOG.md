# Changelog

## [1.4.0](https://github.com/d-strobel/gowindows/compare/v1.3.0...v1.4.0) (2024-12-17)


### Features

* add dhcp client to gowindows client ([524fc19](https://github.com/d-strobel/gowindows/commit/524fc1950ff661b6b7868a59366cadb2b7ebb744))
* add dhcp package for windows dhcp server functions ([b8f1969](https://github.com/d-strobel/gowindows/commit/b8f196901ee3f9b8e2f48c860ca28b44f1610736))
* add dhcp scope read function ([444f7a7](https://github.com/d-strobel/gowindows/commit/444f7a757d181eb75415835363cf0287312fb26a))
* add dhcp ScopeV4Create function ([6df96ca](https://github.com/d-strobel/gowindows/commit/6df96ca444e4f7d8d1c35b552e117c3f9c1896c9))
* add dhcp ScopeV4Delete function ([79c5416](https://github.com/d-strobel/gowindows/commit/79c5416cf8773beb9a56fb72d054095b0234c10e))
* add dhcp ScopeV4Update function ([af522ae](https://github.com/d-strobel/gowindows/commit/af522ae0a74f4c5a9702f774466dd8ccf9ecaa60))
* add helper function for powershell timespans ([db2170e](https://github.com/d-strobel/gowindows/commit/db2170e5f63cd6457cdadb92f070ddd74ff30fa2))
* add ScopeId field to ScopeV4 object ([64ae60e](https://github.com/d-strobel/gowindows/commit/64ae60e922e5d514ed9b3718a842a86d1de92082))
* implement netip package for typed ip address handling ([ce9a05c](https://github.com/d-strobel/gowindows/commit/ce9a05cc5b85eed874ad4774772f8953acda1785))

## [1.3.0](https://github.com/d-strobel/gowindows/compare/v1.2.0...v1.3.0) (2024-12-12)


### Features

* implement UnmarshalJSON for CimTimeDuration ([8df10d5](https://github.com/d-strobel/gowindows/commit/8df10d5659264908b3d94e278f63cb0af0b398bd))

## [1.2.0](https://github.com/d-strobel/gowindows/compare/v1.1.0...v1.2.0) (2024-12-08)


### Features

* add dns server AAAA-Record functions ([584fa32](https://github.com/d-strobel/gowindows/commit/584fa323066786b3109708ea15bde2a30286d53f))
* add dns server functions for read and list dns zones ([f757fc9](https://github.com/d-strobel/gowindows/commit/f757fc9cb700bd2e7959394517d546090caf2c3b))
* add dnsserver A-Record functions ([c154dea](https://github.com/d-strobel/gowindows/commit/c154dea8f8b3704adb627ac1242bbfcc152541dc))
* add dnsserver CName-Record functions ([ffb7267](https://github.com/d-strobel/gowindows/commit/ffb72671af5d6e654e86bc0b7b59f9deaca57b39))
* add dnsserver PTR-Record functions ([8f55816](https://github.com/d-strobel/gowindows/commit/8f558162aabc52e7e230f627716d30a0b2993c61))
* add error handling to parsing json key-value strings ([a4059a7](https://github.com/d-strobel/gowindows/commit/a4059a7ad8ca681286a782e81ec01c7d2179153f))
* implement parsing helper function for json key-value strings ([93da78d](https://github.com/d-strobel/gowindows/commit/93da78dae2104859171546940be7a1e99fca7686))


### Bug Fixes

* cimclass unmarshaljson does not handle array ([6962ec6](https://github.com/d-strobel/gowindows/commit/6962ec6d7d13463083098036893727417d072df4))

## [1.1.0](https://github.com/d-strobel/gowindows/compare/v1.0.1...v1.1.0) (2024-05-11)


### Features

* add winerror package to handle custom errors ([8cb6171](https://github.com/d-strobel/gowindows/commit/8cb617194d0c50f2d89107d97daf475f7e0834a9))
* integrated custom winerror into local accounts package ([849f4c2](https://github.com/d-strobel/gowindows/commit/849f4c27a05b3c4a411afabe243c9ba43a3fee8f))

## [1.0.1](https://github.com/d-strobel/gowindows/compare/v1.0.0...v1.0.1) (2024-05-06)


### Bug Fixes

* dotnettime does not embed the time.Time methods ([32bcba0](https://github.com/d-strobel/gowindows/commit/32bcba09903d79c722d269a9566513525b66449a))

## 1.0.0 (2024-05-06)


### Features

* Add Close function on top client level ([f319b46](https://github.com/d-strobel/gowindows/commit/f319b4605eaa936922761e036fc6062408fbf7cc))
* add context to ssh execution ([645caab](https://github.com/d-strobel/gowindows/commit/645caabc4fc728ebfede9df16e7ed6739d63a9b0))
* add custom error type for parser ([4798e03](https://github.com/d-strobel/gowindows/commit/4798e034be709c90f562e64c93b971ab13a6ce24))
* add custom errors to windows functions ([5620368](https://github.com/d-strobel/gowindows/commit/56203686d9b6b8035ed002d55c3b675db9e6c3ab))
* add errorhandling to run function ([f7f443e](https://github.com/d-strobel/gowindows/commit/f7f443eda46fef926cc40d49f1c696e80779ef04))
* add function to encode powershell commands into powershell.exe commands ([d89a644](https://github.com/d-strobel/gowindows/commit/d89a644c9c07f49cc6cd61df5696ce23a489737f))
* add interface for connection configuration ([8d969a2](https://github.com/d-strobel/gowindows/commit/8d969a2b172cbdb4deffaa518cd998bcea2500b0))
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
* better error message for ssh authentication ([86ba0f5](https://github.com/d-strobel/gowindows/commit/86ba0f5ce7527ad1cb180a277fee465395109fa1))
* change config and connection behaviour ([3443178](https://github.com/d-strobel/gowindows/commit/344317818d49a6a215923ac6fd3d0bb9b70d815f))
* change config methods to private ([8dbd7c1](https://github.com/d-strobel/gowindows/commit/8dbd7c129b2123cc6b51f8c7a1c1793aad159481))
* change defaults method signature ([300626f](https://github.com/d-strobel/gowindows/commit/300626fedc837eee56ae423dcc0893a839e6bc69))
* change function return values ([8f2b39f](https://github.com/d-strobel/gowindows/commit/8f2b39f68d92c60ee7e9cc666825c8f8cda27959))
* Change package name and refactor ([04a6f86](https://github.com/d-strobel/gowindows/commit/04a6f86462ce350fd5ed66d15294f3995504fa70))
* change the config and connection interface ([0e0c17b](https://github.com/d-strobel/gowindows/commit/0e0c17b4c75b8a3b77bea6bf7082297e20f5ab66))
* change WinRMInsecure default value ([5815b4d](https://github.com/d-strobel/gowindows/commit/5815b4d54a6b0e20217553cf41aa782d00f48aba))
* compress json output ([4261f6b](https://github.com/d-strobel/gowindows/commit/4261f6bd6ab2cc8d128e648bd369810205407ed6))
* compress json output ([8515df8](https://github.com/d-strobel/gowindows/commit/8515df8cc34b8d5747eff82602d84997a72f755c))
* delete image ([7863b19](https://github.com/d-strobel/gowindows/commit/7863b190c93834b4733ef69e2fd42e3a441e4723))
* enhanced output for powershell errors ([396c8c5](https://github.com/d-strobel/gowindows/commit/396c8c596ec2609f840fd3324bf9d1ce485a4790))
* implement custom errors ([ac1ccf1](https://github.com/d-strobel/gowindows/commit/ac1ccf1d870826e3078e5fb0e2cd21e6fd46d6aa))
* remove kerberos for now ([e07ed7f](https://github.com/d-strobel/gowindows/commit/e07ed7fd5494ae73b736c8e1678b096eb711bf88))
* remove unneccessary error return ([3ac18f8](https://github.com/d-strobel/gowindows/commit/3ac18f87521f99549acfb3269b79f2dde6edf4aa))
* replace error with custom error ([a26778e](https://github.com/d-strobel/gowindows/commit/a26778eddd3091972aee39656ba36e6c9774635f))
* replace errors with custom errors ([ecd0494](https://github.com/d-strobel/gowindows/commit/ecd049429c326e90eee80efdf389906e6001cb6c))
* replace errors with custom errors ([3466dac](https://github.com/d-strobel/gowindows/commit/3466dac0158e5057ab27bbdf5cf295864d0a4a2f))
* replace errors with custom errors ([c59fa92](https://github.com/d-strobel/gowindows/commit/c59fa92148cf0716f3967e9cb3f49628c9c71bf3))
* switch clients to new connection interface ([ee5e1de](https://github.com/d-strobel/gowindows/commit/ee5e1de233871c02b4dcc9a6cefa7631de0c8da9))
* switch to new powershell encoding function ([73a0dab](https://github.com/d-strobel/gowindows/commit/73a0dab3b2f55ee4ef85d2da57f8a588637f5a70))
* update mocks to new interfaces ([ec1d7a3](https://github.com/d-strobel/gowindows/commit/ec1d7a3571cb08c40a292310dddbedf2b7cfecd9))
* winrm default variables depends on tls ([d54347c](https://github.com/d-strobel/gowindows/commit/d54347c812bdc8ea7e7e13b14af76668d4e073a8))


### Bug Fixes

* cant read group when name has spaces ([d7be6bc](https://github.com/d-strobel/gowindows/commit/d7be6bc2ba9d13c63de054172d843765f90a8318))
* clixml normalizing ([97a01c0](https://github.com/d-strobel/gowindows/commit/97a01c06eceedbfe7fe45342f2337b26265e40d0))
* default ssh known host path ([1dc9c69](https://github.com/d-strobel/gowindows/commit/1dc9c69339a9a8ebacc9150372999fa2b551ff38))
* description param cant be removed with update ([4ae5be7](https://github.com/d-strobel/gowindows/commit/4ae5be70166bcb01990e8f29bfbc6035d616d480))
* gowindows client ([f284073](https://github.com/d-strobel/gowindows/commit/f28407355c042a2d80b0f41217ed05d65fa06e17))
* group name to wildcard leads to error ([9656e61](https://github.com/d-strobel/gowindows/commit/9656e6102a48658da6b3f2798e916c95b7bb868d))
* kerberos config manipulates global winrm TransportDecorator ([e21be43](https://github.com/d-strobel/gowindows/commit/e21be43ec91751e35fcc25f703a1a2275b7e6042))
* local.GroupUpdate cannot update empty description ([23604de](https://github.com/d-strobel/gowindows/commit/23604de252b90524e8086a02588f98a4be673509))
* local.UserCreate: new user without password failes ([f88e252](https://github.com/d-strobel/gowindows/commit/f88e2525282cc25b322fd1e6c00f04a1c02da10d))
* local.UserUpdate function ([38d7580](https://github.com/d-strobel/gowindows/commit/38d7580d57de48d5844e82e3e1dfb5151a7b817c))
* ssh error handling ([c40c785](https://github.com/d-strobel/gowindows/commit/c40c785e2ff42dda7be84d195aab3d729968ffc2))
* ssh returns error if stdout and stderr are empty ([5785490](https://github.com/d-strobel/gowindows/commit/57854903dd6b93c11e2dade6908fe428f0b070a4))
* ssh stdout and stderr are empty ([c5452c0](https://github.com/d-strobel/gowindows/commit/c5452c08b9fc1bf2e5ad8b38f52f38779f8ed8f8))
* stderr handling ([6a4aafb](https://github.com/d-strobel/gowindows/commit/6a4aafb772e4ffcf751f271d6681928ae5b973b2))


### Reverts

* go build version 1.19 and 1.20 ([e5b9187](https://github.com/d-strobel/gowindows/commit/e5b9187508244cf5299cf1e4647cc9cb97d4ccb3))
* move commitlint into build action ([917557e](https://github.com/d-strobel/gowindows/commit/917557e8d7595ca192c6e33ae3b1d44bffd37f3d))
