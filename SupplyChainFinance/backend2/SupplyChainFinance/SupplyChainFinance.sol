// SPDX-License-Identifier: MIT
pragma solidity ^0.6.10;
pragma experimental ABIEncoderV2;

/**
 * @title SupplyChainFinance
 * @dev 供应链金融信息登记系统 - 纯存证，无资金转移
 */
contract SupplyChainFinance {

    // ==================== 核心结构体（3个）====================

    // 1. 企业信息
    struct Enterprise {
        string enterpriseId;        // 企业ID
        string name;                // 企业名称
        address wallet;             // 钱包地址
        bool isCore;                // 是否核心企业
        bool isFinancial;           // 是否金融机构
        uint256 registerTime;       // 注册时间
    }

    // 2. 应收账款
    struct Receivable {
        string receivableId;        // 账单ID
        string orderId;             // 关联订单ID
        address issuer;             // 发行方（供应商）
        address payer;              // 付款方（核心企业）
        uint256 amount;             // 金额
        uint256 issueTime;          // 发行时间
        uint256 dueTime;            // 到期时间
        uint8 status;               // 0-待确认 1-已确认 2-已融资 3-已结算
    }

    // 3. 融资记录
    struct FinancingRecord {
        address financialInst;      // 金融机构
        uint256 amount;             // 融资金额
        uint256 interestRate;       // 利率（万分之几）
        uint256 financingTime;      // 融资时间
    }

    // ==================== 状态常量 ====================
    uint8 constant STATUS_PENDING = 0;
    uint8 constant STATUS_CONFIRMED = 1;
    uint8 constant STATUS_FINANCED = 2;
    uint8 constant STATUS_SETTLED = 3;

    // ==================== 存储映射 ====================

    mapping(address => Enterprise) public enterprises;
    mapping(string => address) public enterpriseIdToAddr;
    address[] public enterpriseList;

    mapping(string => Receivable) public receivables;
    string[] public receivableIdList;

    // 应收账款ID => 融资记录数组
    mapping(string => FinancingRecord[]) public financingRecords;

    // ==================== 事件 ====================

    event EnterpriseRegistered(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial);
    event ReceivableIssued(string receivableId, string orderId, address issuer, address payer, uint256 amount);
    event ReceivableConfirmed(string receivableId, address payer);
    event ReceivableFinanced(string receivableId, address financialInst, uint256 amount, uint256 interestRate);
    event ReceivableSettled(string receivableId, address payer);

    // ==================== 企业注册 ====================

    function registerEnterprise(
        string memory _enterpriseId,
        string memory _name,
        bool _isCore,
        bool _isFinancial
    ) public returns (bool) {
        require(bytes(enterprises[msg.sender].enterpriseId).length == 0, "Already registered");

        enterprises[msg.sender] = Enterprise({
            enterpriseId: _enterpriseId,
            name: _name,
            wallet: msg.sender,
            isCore: _isCore,
            isFinancial: _isFinancial,
            registerTime: block.timestamp
        });

        enterpriseIdToAddr[_enterpriseId] = msg.sender;
        enterpriseList.push(msg.sender);

        emit EnterpriseRegistered(_enterpriseId, _name, msg.sender, _isCore, _isFinancial);
        return true;
    }

    // ==================== 应收账款登记 ====================

    function issueReceivable(
        string memory _receivableId,
        string memory _orderId,
        address _payer,
        uint256 _amount,
        uint256 _dueTime
    ) public returns (bool) {
        require(bytes(receivables[_receivableId].receivableId).length == 0, "Receivable exists");

        receivables[_receivableId] = Receivable({
            receivableId: _receivableId,
            orderId: _orderId,
            issuer: msg.sender,
            payer: _payer,
            amount: _amount,
            issueTime: block.timestamp,
            dueTime: _dueTime,
            status: STATUS_PENDING
        });

        receivableIdList.push(_receivableId);

        emit ReceivableIssued(_receivableId, _orderId, msg.sender, _payer, _amount);
        return true;
    }

    function confirmReceivable(string memory _receivableId) public returns (bool) {
        Receivable storage r = receivables[_receivableId];
        require(bytes(r.receivableId).length > 0, "Not exists");
        require(r.status == STATUS_PENDING, "Not pending");

        r.status = STATUS_CONFIRMED;
        emit ReceivableConfirmed(_receivableId, msg.sender);
        return true;
    }

    // ==================== 融资登记（纯记录）====================

    function recordFinancing(
        string memory _receivableId,
        uint256 _amount,
        uint256 _interestRate
    ) public returns (bool) {
        Receivable storage r = receivables[_receivableId];
        require(bytes(r.receivableId).length > 0, "Not exists");
        require(r.status == STATUS_CONFIRMED, "Not confirmed");

        financingRecords[_receivableId].push(FinancingRecord({
            financialInst: msg.sender,
            amount: _amount,
            interestRate: _interestRate,
            financingTime: block.timestamp
        }));

        r.status = STATUS_FINANCED;

        emit ReceivableFinanced(_receivableId, msg.sender, _amount, _interestRate);
        return true;
    }

    // ==================== 结算登记（纯状态变更）====================

    function recordSettlement(string memory _receivableId) public returns (bool) {
        Receivable storage r = receivables[_receivableId];
        require(bytes(r.receivableId).length > 0, "Not exists");
        require(r.status == STATUS_CONFIRMED || r.status == STATUS_FINANCED, "Invalid status");

        r.status = STATUS_SETTLED;
        emit ReceivableSettled(_receivableId, msg.sender);
        return true;
    }

    // ==================== 查询函数 ====================

    function getEnterprise(address _addr) public view returns (Enterprise memory) {
        return enterprises[_addr];
    }

    function getReceivable(string memory _id) public view returns (Receivable memory) {
        return receivables[_id];
    }

    function getFinancingRecords(string memory _receivableId) public view returns (FinancingRecord[] memory) {
        return financingRecords[_receivableId];
    }

    function getEnterpriseCount() public view returns (uint256) {
        return enterpriseList.length;
    }

    function getReceivableCount() public view returns (uint256) {
        return receivableIdList.length;
    }

    function getReceivableIdByIndex(uint256 _index) public view returns (string memory) {
        require(_index < receivableIdList.length, "Out of bounds");
        return receivableIdList[_index];
    }
}