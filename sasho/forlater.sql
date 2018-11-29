CREATE TABLE `Blocks` (
	`UnclesCount`	INTEGER,
	`TxCount`	INTEGER,
	`Number`	INTEGER,
	`GasLimit`	NUMERIC,
	`GasUsed`	NUMERIC,
	`Difficulty`	TEXT,
	`Time`	TEXT,
	`Nonce`	TEXT,
	`Miner`	TEXT,
	`ParentHash`	TEXT,
	`Hash`	TEXT,
	`TxRoot`	TEXT,
	`ExtraData`	TEXT,
	PRIMARY KEY(Number)
);

INSERT INTO Blocks (
    UnclesCount,
	TxCount,
	Number,
	GasLimit,
	GasUsed,
	Difficulty,
	Time,
	Nonce,
	Miner,
	ParentHash,
	Hash,
	ExtraData
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)

CREATE TABLE `Headers` (
	`ParentHash`	TEXT,
	`Sha3Uncles`	TEXT,
	`Miner`	TEXT,
	`StateRoot`	TEXT,
	`TransactionsRoot`	TEXT,
	`ReceiptsRoot`	TEXT,
	`LogsBloom`	TEXT,
	`Difficulty`	NUMERIC,
	`Number`	NUMERIC,
	`GasLimit`	NUMERIC,
	`GasUsed`	NUMERIC,
	`Timestamp`	TEXT,
	`ExtraData`	TEXT,
	`MixHash`	TEXT,
	`Nonce`	NUMERIC,
    `Hash` TEXT
);