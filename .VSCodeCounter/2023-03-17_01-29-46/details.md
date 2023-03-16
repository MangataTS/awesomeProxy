# Details

Date : 2023-03-17 01:29:46

Directory e:\\goland\\project\\awesomeProxy

Total : 50 files,  4838 codes, 936 comments, 728 blanks, all 6502 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [.idea/awesomeProxy.iml](/.idea/awesomeProxy.iml) | XML | 9 | 0 | 0 | 9 |
| [.idea/modules.xml](/.idea/modules.xml) | XML | 8 | 0 | 0 | 8 |
| [.idea/vcs.xml](/.idea/vcs.xml) | XML | 6 | 0 | 0 | 6 |
| [Cert/WindowsCA.go](/Cert/WindowsCA.go) | Go | 49 | 8 | 11 | 68 |
| [Constant/Certificate.go](/Constant/Certificate.go) | Go | 1 | 0 | 2 | 3 |
| [Contract/IServerProcesser.go](/Contract/IServerProcesser.go) | Go | 4 | 0 | 4 | 8 |
| [Core/Certificate.go](/Core/Certificate.go) | Go | 168 | 10 | 10 | 188 |
| [Core/ConnPeer.go](/Core/ConnPeer.go) | Go | 11 | 0 | 3 | 14 |
| [Core/ProxyHttp.go](/Core/ProxyHttp.go) | Go | 428 | 42 | 21 | 491 |
| [Core/ProxyServer.go](/Core/ProxyServer.go) | Go | 99 | 5 | 12 | 116 |
| [Core/ProxySocket5.go](/Core/ProxySocket5.go) | Go | 246 | 21 | 11 | 278 |
| [Core/ProxyTcp.go](/Core/ProxyTcp.go) | Go | 91 | 2 | 6 | 99 |
| [Core/Storage.go](/Core/Storage.go) | Go | 67 | 4 | 9 | 80 |
| [Core/Websocket/Client.go](/Core/Websocket/Client.go) | Go | 271 | 74 | 51 | 396 |
| [Core/Websocket/ClientClone.go](/Core/Websocket/ClientClone.go) | Go | 8 | 4 | 5 | 17 |
| [Core/Websocket/ClientCloneLegacy.go](/Core/Websocket/ClientCloneLegacy.go) | Go | 26 | 8 | 5 | 39 |
| [Core/Websocket/Compression.go](/Core/Websocket/Compression.go) | Go | 117 | 11 | 21 | 149 |
| [Core/Websocket/Conn.go](/Core/Websocket/Conn.go) | Go | 876 | 171 | 155 | 1,202 |
| [Core/Websocket/ConnWrite.go](/Core/Websocket/ConnWrite.go) | Go | 7 | 4 | 5 | 16 |
| [Core/Websocket/ConnWriteLegacy.go](/Core/Websocket/ConnWriteLegacy.go) | Go | 11 | 4 | 4 | 19 |
| [Core/Websocket/Doc.go](/Core/Websocket/Doc.go) | Go | 1 | 225 | 2 | 228 |
| [Core/Websocket/Join.go](/Core/Websocket/Join.go) | Go | 31 | 6 | 6 | 43 |
| [Core/Websocket/Json.go](/Core/Websocket/Json.go) | Go | 34 | 20 | 7 | 61 |
| [Core/Websocket/Mask.go](/Core/Websocket/Mask.go) | Go | 35 | 9 | 11 | 55 |
| [Core/Websocket/MaskSafe.go](/Core/Websocket/MaskSafe.go) | Go | 8 | 4 | 4 | 16 |
| [Core/Websocket/Prepared.go](/Core/Websocket/Prepared.go) | Go | 70 | 20 | 13 | 103 |
| [Core/Websocket/Proxy.go](/Core/Websocket/Proxy.go) | Go | 60 | 5 | 13 | 78 |
| [Core/Websocket/Server.go](/Core/Websocket/Server.go) | Go | 213 | 111 | 44 | 368 |
| [Core/Websocket/Trace.go](/Core/Websocket/Trace.go) | Go | 15 | 1 | 4 | 20 |
| [Core/Websocket/Trace_17.go](/Core/Websocket/Trace_17.go) | Go | 8 | 1 | 4 | 13 |
| [Core/Websocket/Util.go](/Core/Websocket/Util.go) | Go | 244 | 26 | 14 | 284 |
| [Core/Websocket/XnetProxy.go](/Core/Websocket/XnetProxy.go) | Go | 352 | 57 | 65 | 474 |
| [Log/Logger.go](/Log/Logger.go) | Go | 549 | 46 | 61 | 656 |
| [README.md](/README.md) | Markdown | 81 | 0 | 28 | 109 |
| [Reproxy/handle.go](/Reproxy/handle.go) | Go | 76 | 6 | 13 | 95 |
| [Utils/Utils.go](/Utils/Utils.go) | Go | 26 | 0 | 4 | 30 |
| [balance/balance_mgr.go](/balance/balance_mgr.go) | Go | 24 | 0 | 7 | 31 |
| [balance/blancer.go](/balance/blancer.go) | Go | 4 | 0 | 2 | 6 |
| [balance/hash_balance.go](/balance/hash_balance.go) | Go | 27 | 1 | 8 | 36 |
| [balance/instance.go](/balance/instance.go) | Go | 29 | 0 | 8 | 37 |
| [balance/random_balance.go](/balance/random_balance.go) | Go | 20 | 1 | 7 | 28 |
| [balance/round_robin_balance.go](/balance/round_robin_balance.go) | Go | 23 | 1 | 9 | 33 |
| [balance/shuffle2_balance.go](/balance/shuffle2_balance.go) | Go | 26 | 2 | 9 | 37 |
| [balance/shuffle_balance.go](/balance/shuffle_balance.go) | Go | 26 | 2 | 9 | 37 |
| [balance/weight_round_robin.go](/balance/weight_round_robin.go) | Go | 65 | 4 | 15 | 84 |
| [config.json](/config.json) | JSON | 22 | 0 | 0 | 22 |
| [config/ConfigInit.go](/config/ConfigInit.go) | Go | 169 | 8 | 14 | 191 |
| [go.mod](/go.mod) | Go Module File | 6 | 0 | 3 | 9 |
| [go.sum](/go.sum) | Go Checksum File | 4 | 0 | 1 | 5 |
| [main.go](/main.go) | Go | 87 | 12 | 8 | 107 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)