# bstudent_jc

## 1. Network 실행

network 경로에서 (다음 쉘스크립트 순차 실행):

- ./startJongchinNet.sh 
- ./deployJongchinCC.sh

## 2. CCP 생성

application/ccp 경로에서 (다음 쉘스크립트 실행):

- ./ccp-generate.sh

## 3. 인증서 가져오기 및 지갑 생성

application 경로에서 (다음 쉘스크립트 순차 실행):

- ./getCert.sh
- ./addToWallet.js

## 4. node.js 모듈 설치

application 경로에서

- npm install
