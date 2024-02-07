# Overview

This document describe the main idea of the Merkle Tree challange, who it's implemented and how to run/show.
The main requirements:
- User wants to has an ability to upload files (which are not too big in size)
- After uploading user does not want to keep them locally
- Also, user wants to have an ability to download the particular file and make sure that it's actually the same file as was uploaded.
So here we need to ensure somehow that file was not corrupted by transport issues or vulnerabilities.

# Approach

For make it happens we may use Merkle Tree data structure. It's a binary tree where each non-leaf node is a hash of its children. The root is a hash of all the blocks. So, if we have a hash of the root, we can be sure that all the blocks are the same as they were when we calculated the root hash. Also, we can check the integrity of the particular block by checking the hash of the root and the hashes of the blocks on the path from the root to the block.
List of hashes calculates for list of files one by one and then we can calculate the root hash of the tree. On client side we can store the root hash and calculate it each time when the server sends us the file with Merkle Proof. The Merkle Proof is a list of hashes from the root to the block. We can check the integrity of the block by comparing the hash of the root and the hashes from the Merkle Proof. If the calculated root hash is the same as we have stored, we can be sure that all the files are the same as they were when we calculated the root hash.

# Implementation

The server is implemented using `echo` framework. It has 2 endpoints:
- `POST /upload` - for uploading files
- `GET /download` - for downloading files

The server stores files in file system and calculates the Merkle Tree when the file should be downloaded. The server sends the file with the Merkle Proof to the client. It utilizes `mime/multipart` content-type and includes 2 fields:
- `file` - the file itself. (Bytes)
- `proof` - the Merkle Proof. (JSON). The proof is a list of hashes and indices of the blocks in the tree.

The client is implemented using `cli` framework. It has 3 commands:
- `upload` - for uploading files
- `download` - for downloading files
- `generate` - for generating random files

`download` command has an additional flag `--corrupt` which allows to corrupt the file by changing bytes of received file. It uses only for demo purposes to show that if content was changed somehow file will not be the same as it was when we calculated the root hash and will not be written to the file system.

## Proof

### RootHash

It's needed only for the client side for making the first check. If it's not the same as client stored we can be sure that the file was changed and we don't need to check the Merkle Proof at all.

### Hashes

Hashes is a list of hashes from the root to the block. The first hash is the hash of the block, the last hash is the hash of the root. For example, if we have the following tree:
```
      root
     /    \
    a      b
   / \    / \
  c   d  e   f
```
hashes for the block `e` will be:
```
"hashes": "hash(f), hash(a)"
```
We don't need to send hash(e) cause it will be calculated on the client side from the File content. Also client store rootHash and can calculate it each time when the server sends the file with Merkle Proof.

### Indices

Indices are indicates the directioin of the hash calculation. "0" means that hash came from the left child, "1" means that hash came from the right child.
For example, if we have the following tree:
```
      root
     /    \
    a      b
   / \    / \
  c   d  e   f
```

The proof for the block `e` will be:
```
"proof": {
    "hashes": "hash(f), hash(a)",     
    "indices": "1, 0"
}
```

# Future improvements

- Store calculated MerkleTree on server side in database. It will allow to avoid calculating the tree each time when the file should be downloaded.
- Make file downloading more async. It will allow to download file with huge size.
- Move `indices` from list to bitmask. It will reduce the size of the proof.
- Add integration tests based on docker-compose. For now it's used only for demo purposes.

# What went well/bad

What was good:
- The main implementation of server-client.
- Merkle Tree building/calculation.
- Dockerizing
- Demo

What was not good and took more time:
- `mime/multipart` content-type. It's good documented for `echo` framework but spent some time on understanding how it handle on the client side (resty).
- Idea to use indices for the proof. Spent some time on investigating how it's implemented on the real projects. (Was insipred by https://github.com/cbergoon/merkletree project)
- Spent some time on cobra framework.

# Demo

Demo is stored in `docker` folder and utilizes `docker-compose` for running the server and the client. It runs `demo-client.sh` script under the hood.

## Demo steps

1. Generate random files
2. Upload files
3. Download file Positive case
4. Download file Negative case (--corrupt option which corrupts the bytes of received file on purpose)

## Preps

Go to `docker` folder and run the following command:
```bash
    docker compose build
```

## Run

For running the demo use the following command:
```bash
    docker compose up
```

Example of the output:
```text
~ docker compose up
[+] Running 3/3
 ⠿ Network docker_default     Created                                                                                                                                                                 0.1s
 ⠿ Container docker-server-1  Created                                                                                                                                                                 0.1s
 ⠿ Container docker-client-1  Created                                                                                                                                                                 0.0s
Attaching to docker-client-1, docker-server-1
docker-server-1  | 1. Create a directory for storage
docker-server-1  | 2. Start the server
docker-server-1  | 
docker-server-1  |    ____    __
docker-server-1  |   / __/___/ /  ___
docker-server-1  |  / _// __/ _ \/ _ \
docker-server-1  | /___/\__/_//_/\___/ v4.11.4
docker-server-1  | High performance, minimalist Go web framework
docker-server-1  | https://echo.labstack.com
docker-server-1  | ____________________________________O/_______
docker-server-1  |                                     O\
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:31Z","message":"Starting listener"}
docker-server-1  | ⇨ http server started on [::]:8080
docker-client-1  | 0. Create a directory for storage
docker-client-1  | 1. Generate let's say 10 files with random content, by calling ./client generate 10
docker-client-1  | File in storage: file0
docker-client-1  | File in storage: file1
docker-client-1  | File in storage: file2
docker-client-1  | File in storage: file3
docker-client-1  | File in storage: file4
docker-client-1  | File in storage: file5
docker-client-1  | File in storage: file6
docker-client-1  | File in storage: file7
docker-client-1  | File in storage: file8
docker-client-1  | File in storage: file9
docker-client-1  | 2. Upload that files to server, by calling "./client upload"
docker-client-1  | {"level":"info","time":"2024-02-07T13:34:31Z","message":"Merkle Root Hash: aa090ca39c0ae121a86981dc7f07f7d2fef6491022390300eaed02b66f80fc20"}
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"Uploading files..."}
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"Parsed: &{map[] map[files:[0xc000054480 0xc0000544e0 0xc000054540 0xc0000545a0 0xc000054600 0xc000054660 0xc0000546c0 0xc000054720 0xc000054780 0xc0000547e0]]}"}
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"Files: [0xc000054480 0xc0000544e0 0xc000054540 0xc0000545a0 0xc000054600 0xc000054660 0xc0000546c0 0xc000054720 0xc000054780 0xc0000547e0]"}
docker-server-1  | {"time":"2024-02-07T13:34:32.004752469Z","id":"","remote_ip":"192.168.144.1","host":"0.0.0.0:8080","method":"POST","uri":"/files","user_agent":"go-resty/2.10.0 (https://github.com/go-resty/resty)","status":200,"error":"","latency":2673566,"latency_human":"2.673566ms","bytes_in":1836,"bytes_out":32}
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"Merkle tree generated: aa090ca39c0ae121a86981dc7f07f7d2fef6491022390300eaed02b66f80fc20"}
docker-client-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"Upload response: Files were uploaded successfully"}
docker-client-1  | 3. Check, that files were removed from storage.
docker-client-1  | File in storage: root_hash
docker-client-1  | # 4. Check that only file `root_hash` is present in storage.
docker-client-1  | # 5. Download file from server, by calling "./client download file1"
docker-client-1  | # 6. Check that file was added to storage.
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"File file1 exists: true"}
docker-server-1  | {"time":"2024-02-07T13:34:32.016806134Z","id":"","remote_ip":"192.168.144.1","host":"0.0.0.0:8080","method":"GET","uri":"/files/file1","user_agent":"go-resty/2.10.0 (https://github.com/go-resty/resty)","status":200,"error":"","latency":412336,"latency_human":"412.336µs","bytes_in":0,"bytes_out":853}
docker-client-1  | File in storage: file1
docker-client-1  | File in storage: root_hash
docker-client-1  | # 7. Download file from server with corruption imitation, by calling "./client download file7 --corrupt"
docker-server-1  | {"time":"2024-02-07T13:34:32.035889299Z","id":"","remote_ip":"192.168.144.1","host":"0.0.0.0:8080","method":"GET","uri":"/files/file7","user_agent":"go-resty/2.10.0 (https://github.com/go-resty/resty)","status":200,"error":"","latency":398168,"latency_human":"398.168µs","bytes_in":0,"bytes_out":853}
docker-server-1  | {"level":"info","time":"2024-02-07T13:34:32Z","message":"File file7 exists: true"}
docker-client-1  | Invalid proof. Seems like File is corrupted
docker-client-1  | # 8. Check that file was not added to storage.
docker-client-1  | File in storage: file1
docker-client-1  | File in storage: root_hash
docker-client-1  | Demo is finished
docker-client-1 exited with code 0
```