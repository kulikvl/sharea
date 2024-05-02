Server (Go):
- validate flags: 
  - capacity - if space of disk exists
- api: 
  - better error handling (status codes and messages)
  - automatically rename uploaded file if already exists with the same name
  - make var like "uploadBufferSize" and test what size is optimal
  - check if it is possible to make more efficient download system
- test:
  - unit tests where possible
  - think about mocks for integrations tests
- config:
  - analyse if embedding a Web FS would be beneficial
- watcher:
  - better handle errors and check the code overall
- qr:
  - add qr-code scanning of server address

Client (JS);
- js:
  - refactor and analyse index.js
- html&css:
  - better styles

Overall:
- write README.md
- write .goreleaser.yml