---
date: 2026-05-18
time: 12:47PM
tags: post
snippet: yes
layout: blog.liquid
description: yes
---
You might find this useful if you're doing anything with EDL to unlock/use adb on the Light Phone II. Partition UUIDs will probably be different on your device.
```
Keystone library is missing (optional).
Qualcomm Sahara / Firehose Client V3.62 (c) B.Kerler 2018-2025.
main - Trying with no loader given ...
main - Waiting for the device
main - Device detected :)
sahara - Protocol version: 2, Version supported: 1
main - Mode detected: sahara
sahara -
Version 0x2
------------------------
HWID:              0x009600e100000000 (MSM_ID:0x009600e1,OEM_ID:0x0000,MODEL_ID:0x0000)
CPU detected:      "MSM8909"
PK_HASH:           0xcc3153a80293939b90d02d3bf8b23e0292e452fef662c74998421adad42a380f
Serial:            0x23344139

sahara - Possibly unfused device detected, so any loader should be fine...
sahara - Trying loader: ./Loaders/qualcomm/model_generic/msm8909/009600e100000000_cc3153a80293939b_fhprg_peek.bin
sahara - Protocol version: 2, Version supported: 1
sahara - Uploading loader ./Loaders/qualcomm/model_generic/msm8909/009600e100000000_cc3153a80293939b_fhprg_peek.bin ...
sahara - 32-Bit mode detected.
sahara - Loader successfully uploaded.
main - Trying to connect to firehose loader ...
firehose_client
firehose_client - [LIB]: No --memory option set, we assume "eMMC" as default ..., if it fails, try using "--memory" with "UFS","NAND" or "spinor" instead !
firehose - TargetName=MSM8909
firehose - MemoryName=eMMC
firehose - Version=1
firehose - Trying to read first storage sector...
firehose - Running configure...
firehose_client - Supported functions:
-----------------

Parsing Lun 0:

GPT Table:
-------------
modem:               Offset 0x0000000004000000, Length 0x0000000004200000, Flags 0x1000000000000000, UUID aafbc0d1-c6bc-074a-f815-2df473662a68, Type EFI_BASIC_DATA, Active False
DDR:                 Offset 0x0000000008200000, Length 0x0000000000008000, Flags 0x1000000000000000, UUID 813402be-c6b1-5854-0646-30ec5d049e5d, Type 0x20a0c19c, Active False
fsg:                 Offset 0x0000000008208000, Length 0x0000000000180000, Flags 0x1000000000000000, UUID 288c0614-9fec-bdc0-aaa4-4b7548006fee, Type 0x638ff8e2, Active False
sec:                 Offset 0x0000000008388000, Length 0x0000000000004000, Flags 0x1000000000000000, UUID 92b3411a-b431-5d5f-c5f1-21455418cddf, Type 0x303e6ac3, Active False
boot:                Offset 0x000000000838c000, Length 0x0000000002000000, Flags 0x1000000000000000, UUID 27a93195-4af4-d9a5-ea58-1aac4954a27f, Type 0x20117f86, Active False
system:              Offset 0x000000000a38c000, Length 0x000000004b000000, Flags 0x1000000000000000, UUID 121e129c-1604-9d25-95ec-c7d3bd0f727a, Type 0x97d7b011, Active False
persist:             Offset 0x000000005538c000, Length 0x0000000002000000, Flags 0x1000000000000000, UUID 02b9ca8c-88f1-66a5-6911-793a36618b1b, Type 0x6c95e238, Active False
cache:               Offset 0x000000005738c000, Length 0x0000000006e00000, Flags 0x1000000000000000, UUID 357b9725-d8c5-60cc-e95a-9c814ed8c50d, Type 0x5594c694, Active False
recovery:            Offset 0x000000005e18c000, Length 0x0000000002000000, Flags 0x1000000000000000, UUID a1a8f83d-7c81-3846-1b12-70dc1dbbdeef, Type 0x9d72d4e4, Active False
devinfo:             Offset 0x000000006018c000, Length 0x0000000000100000, Flags 0x1000000000000000, UUID 4fb06217-c5bf-4609-353a-5d4b9a22872f, Type 0x1b81e7e6, Active False
cmnlib:              Offset 0x000000006028c000, Length 0x0000000000040000, Flags 0x1000000000000000, UUID 1a5412e5-2ee2-8cf0-40b0-da3ed11e5e1a, Type 0x73471795, Active False
cmnlibbak:           Offset 0x00000000602cc000, Length 0x0000000000040000, Flags 0x1000000000000000, UUID 404c6206-2b6d-5f6d-b7c0-85680401d1a1, Type 0x73471795, Active False
keymaster:           Offset 0x000000006030c000, Length 0x0000000000080000, Flags 0x1000000000000000, UUID 1b94a3db-de9a-974a-23a7-311f1a19d2f6, Type 0xe8b7cf6e, Active False
keymasterbak:        Offset 0x000000006038c000, Length 0x0000000000080000, Flags 0x1000000000000000, UUID 9d02327f-fe4c-09b2-9127-80ca87795102, Type 0xe8b7cf6e, Active False
sbl1:                Offset 0x0000000064000000, Length 0x0000000000080000, Flags 0x0000000000000000, UUID 8bff55b8-9ab4-36a4-8e24-58879d1c72f0, Type 0xdea0ba2c, Active False
sbl1bak:             Offset 0x0000000064080000, Length 0x0000000000080000, Flags 0x0000000000000000, UUID b1baaf15-4bf4-3885-1817-1af2e740aeb9, Type 0xdea0ba2c, Active False
aboot:               Offset 0x0000000064100000, Length 0x0000000000100000, Flags 0x0000000000000000, UUID ffb79caf-3476-3b3d-7722-8a9eaa4feae9, Type 0x400ffdcd, Active False
abootbak:            Offset 0x0000000064200000, Length 0x0000000000100000, Flags 0x0000000000000000, UUID 52781cac-fd14-ebcc-60a9-93076503817a, Type 0x400ffdcd, Active False
rpm:                 Offset 0x0000000064300000, Length 0x0000000000080000, Flags 0x0000000000000000, UUID 4d7ec450-0315-0371-5a23-fb14c2ce56b0, Type 0x98df793, Active False
rpmbak:              Offset 0x0000000064380000, Length 0x0000000000080000, Flags 0x0000000000000000, UUID 897bfe74-df22-884d-4bbd-8123fe15dc40, Type 0x98df793, Active False
tz:                  Offset 0x0000000064400000, Length 0x0000000000200000, Flags 0x0000000000000000, UUID 2fe2ae56-4fcc-79eb-6686-91dd55a99477, Type 0xa053aa7f, Active False
tzbak:               Offset 0x0000000064600000, Length 0x0000000000200000, Flags 0x0000000000000000, UUID 968b570e-242c-3d0d-0880-a6ed5232c8ac, Type 0xa053aa7f, Active False
devcfg:              Offset 0x0000000064800000, Length 0x0000000000040000, Flags 0x0000000000000000, UUID 64e567b2-b883-4337-b595-97851fc41d78, Type 0xf65d4b16, Active False
apdp:                Offset 0x0000000064840000, Length 0x0000000000040000, Flags 0x0000000000000000, UUID fb8de2c9-9744-4080-3b60-9fc243e5955f, Type 0xe6e98da2, Active False
pad:                 Offset 0x0000000064880000, Length 0x0000000000100000, Flags 0x0000000000000000, UUID cacd76e0-752b-f39d-819d-cafcdc3a0a36, Type EFI_BASIC_DATA, Active False
modemst1:            Offset 0x0000000064980000, Length 0x0000000000180000, Flags 0x0000000000000000, UUID 8b9a6b28-4c4f-f81f-3bd8-234a777e8c4f, Type 0xebbeadaf, Active False
modemst2:            Offset 0x0000000064b00000, Length 0x0000000000180000, Flags 0x0000000000000000, UUID d2ab7078-07d1-ccc4-6bc5-d84fb46894e0, Type 0xa288b1f, Active False
misc:                Offset 0x0000000064c80000, Length 0x0000000000800000, Flags 0x0000000000000000, UUID 09b96c23-94b8-5acd-1a8d-663def3ab93c, Type 0x82acc91f, Active False
fsc:                 Offset 0x0000000065480000, Length 0x0000000000000400, Flags 0x0000000000000000, UUID 4dac0f31-b7c4-f0bb-c2b7-521d44de13be, Type 0x57b90a16, Active False
ssd:                 Offset 0x0000000065480400, Length 0x0000000000002000, Flags 0x0000000000000000, UUID 8c6835e3-bf4d-cd7a-ce33-3524ea59335b, Type 0x2c86e742, Active False
splash:              Offset 0x0000000065482400, Length 0x0000000000a00000, Flags 0x0000000000000000, UUID 1a6002c4-9b1d-233e-8669-3e9d33d64925, Type 0x20117f86, Active False
keystore:            Offset 0x0000000065e82400, Length 0x0000000000080000, Flags 0x0000000000000000, UUID 69b98fee-40eb-e8a2-367b-95ca0f426f46, Type 0xde7d4029, Active False
oem:                 Offset 0x0000000065f02400, Length 0x0000000000800000, Flags 0x0000000000000000, UUID 145912d4-e78f-d3ad-6e6a-c6639f93fa63, Type 0x7db6ac55, Active False
prodinfo:            Offset 0x0000000066702400, Length 0x0000000000001000, Flags 0x0000000000000000, UUID f65166b1-69bb-5716-cde7-4a4df7d0de9f, Type 0x21130059, Active False
prodinfo2:           Offset 0x0000000066703400, Length 0x0000000000001000, Flags 0x0000000000000000, UUID 08a31a81-fa15-109f-a84c-a87936df0ea4, Type 0xbc3c23ce, Active False
config:              Offset 0x0000000066704400, Length 0x0000000000080000, Flags 0x0000000000000000, UUID 0639975e-1be7-87a9-472f-079da2fd561f, Type 0x91b72d4d, Active False
vendor:              Offset 0x0000000066784400, Length 0x000000000fa00000, Flags 0x0000000000000000, UUID e0823c69-76f4-008a-5e66-db089530337f, Type 0x3c160a98, Active False
userdata:            Offset 0x0000000076184400, Length 0x000000015be77a00, Flags 0x0000000000000000, UUID 03ea9a05-1ec1-bad5-7bad-9c3949da6ca8, Type 0x1b81e7e6, Active False

Total disk size:0x00000001d2000000, sectors:0x0000000000e90000
```