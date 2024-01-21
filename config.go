package dbgo

// 'mysql' => [
//    'read' => [
//        'host' => [
//            '192.168.1.1',
//            '196.168.1.2',
//        ],
//    ],
//    'write' => [
//        'host' => [
//            '196.168.1.3',
//        ],
//    ],
//    'sticky' => true,
//    'driver' => 'mysql',
//    'database' => 'database',
//    'username' => 'root',
//    'password' => '',
//    'charset' => 'utf8mb4',
//    'collation' => 'utf8mb4_unicode_ci',
//    'prefix' => '',
//],

type Config struct {
	Host      []string
	Sticky    bool   `toml:"sticky"`
	Driver    string `toml:"driver"`
	Database  string `toml:"database"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Charset   string `toml:"charset"`
	Collation string `toml:"collation"`
	Prefix    string `toml:"prefix"`
}

type Cluster struct {
	Read      []string
	Write     []string
	Sticky    bool   `toml:"sticky"`
	Driver    string `toml:"driver"`
	Database  string `toml:"database"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Charset   string `toml:"charset"`
	Collation string `toml:"collation"`
	Prefix    string `toml:"prefix"`
}
