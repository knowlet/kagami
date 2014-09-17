CREATE TABLE `accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` char(12) NOT NULL,
  `password` char(128) NOT NULL,
  `salt` char(10) DEFAULT NULL,
  `char_delete_password` int(8) unsigned NOT NULL,
  `online` tinyint(1) NOT NULL DEFAULT '0',
  `banned` tinyint(1) NOT NULL DEFAULT '0',
  `ban_expire` datetime DEFAULT NULL,
  `ban_reason` tinyint(2) unsigned DEFAULT NULL,
  `ban_reason_message` varchar(255) DEFAULT NULL,
  `last_login` datetime DEFAULT NULL,
  `creation_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username_UNIQUE` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ip_bans` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_UNIQUE` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;