/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

-- This database is heavy based on MapleStory Vana

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
  `creation_date` datetime DEFAULT NULL,
  `admin` tinyint(1) NOT NULL DEFAULT '0',
  `gm_level` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username_UNIQUE` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ip_bans` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_UNIQUE` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `characters` (
  `character_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(12) NOT NULL,
  `user_id` int(11) NOT NULL,
  `world_id` tinyint(3) unsigned NOT NULL,
  `level` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `job` smallint(6) NOT NULL DEFAULT '0',
  `str` smallint(6) NOT NULL DEFAULT '4',
  `dex` smallint(6) NOT NULL DEFAULT '4',
  `int` smallint(6) NOT NULL DEFAULT '4',
  `luk` smallint(6) NOT NULL DEFAULT '4',
  `chp` smallint(6) NOT NULL DEFAULT '50',
  `mhp` smallint(6) NOT NULL DEFAULT '50',
  `cmp` smallint(6) NOT NULL DEFAULT '5',
  `mmp` smallint(6) NOT NULL DEFAULT '5',
  `ap` smallint(6) NOT NULL DEFAULT '9',
  `sp` smallint(6) NOT NULL DEFAULT '0',
  `exp` int(11) NOT NULL DEFAULT '0',
  `fame` smallint(6) NOT NULL DEFAULT '0',
  `map` int(11) NOT NULL DEFAULT '0',
  `pos` smallint(6) NOT NULL DEFAULT '0',
  `gender` tinyint(1) NOT NULL,
  `skin` tinyint(4) NOT NULL,
  `face` int(11) NOT NULL,
  `hair` int(11) NOT NULL,
  `online` tinyint(1) NOT NULL DEFAULT '0',
  `overall_cpos` int(11) unsigned DEFAULT NULL,
  `overall_opos` int(11) unsigned DEFAULT NULL,
  `world_cpos` int(11) unsigned DEFAULT NULL,
  `world_opos` int(11) unsigned DEFAULT NULL,
  `job_cpos` int(11) unsigned DEFAULT NULL,
  `job_opos` int(11) unsigned DEFAULT NULL,
  `fame_cpos` int(11) unsigned DEFAULT NULL,
  `fame_opos` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`character_id`),
  KEY `user_id` (`user_id`),
  KEY `world_id` (`world_id`),
  KEY `name` (`name`),
  CONSTRAINT `characters_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `items` (
  `character_id` int(11) NOT NULL,
  `inv` smallint(6) NOT NULL,
  `slot` smallint(6) NOT NULL,
  `location` enum('inventory','storage') NOT NULL,
  `user_id` int(11) NOT NULL,
  `world_id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  PRIMARY KEY (`character_id`,`inv`,`slot`,`location`),
  CONSTRAINT `items_ibfk_1` FOREIGN KEY (`character_id`) REFERENCES `characters` (`character_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `storage` (
  `user_id` int(11) NOT NULL,
  `world_id` int(11) NOT NULL,
  `slots` smallint(6) NOT NULL,
  `mesos` int(11) NOT NULL,
  `char_slots` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`world_id`),
  CONSTRAINT `storage_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
