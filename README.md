# taskbytime
타임기반의 일처리 엔진

사용처 -----
1. 하트, 날개, 신발등등 오브젝트 생성
2. 시간이 걸리는 건물 건설 / 능력 업그레이드

기능 사양 -----
1. 태스크를 시작하고 인터벌단위로 수량 생성
2. 최대수량 초과되면 생성기능 정지
3. Add API 에 의해 수량 관리
4. 작업1 당 DB열 1

설계 방향 -----
1. 시작 시간등의 정보를 저장해 놓고 API 호출됐을 때만 계산함 -> 최적화
2. 모든 수량, 시간의 정보는 서버에 의해 동기화 됨 (현재수량, 인터벌, remainTime)
3. task index 은 각각 시작수량, 최대수량, 인터벌, 옵션(Reduce 시에 0이 되면 task 삭제)의 파라메터를 가지고 있음

제공 API -----
1. Create
    설명 : 태스크를 생성함
    인풋 : user id, task index
    리턴 : 현재수량, 인터벌, err
1. Start
    설명 : 인터벌 시작. 현재수량이 최대수량보다 크거나 같을 경우에는 에러처리
    인풋 : user id, task index
    리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
2. Add :
    설명 : 수량을 늘림 (친구의 하트선물등)
    인풋 : user id, task index, add number
    리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
3. Reduce :
    설명 : 수량을 줄임. 자동 삭제(옵션)
    인풋 : user id, task index, reduce number
    리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
4. Delete :
    설명 : 태스크 삭제
    인풋 : user id, task index
    리턴 : err
