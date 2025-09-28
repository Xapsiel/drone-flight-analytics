
CREATE TABLE  IF NOT EXISTS messages (
                                 id SERIAL PRIMARY KEY ,
                                 region INTEGER  REFERENCES district_shapes(gid),
                                 sid VARCHAR(100) UNIQUE ,
                                 dof DATE NOT NULL,
                                 atd TIME NOT NULL,
                                 ata TIME,
                                 dep_coords_normalize VARCHAR(20) NOT NULL,
                                 arr_coords_normalize VARCHAR(20),
                                 dep_coordinate geography(POINT, 4326) NOT NULL,
                                 arr_coordinate geography(POINT, 4326),
                                 arr_region_rf VARCHAR(255),
                                 opr VARCHAR(100),
                                 reg VARCHAR(50),
                                 typ VARCHAR(10),
                                 rmk TEXT,
                                 min_alt INTEGER NOT NULL,
                                 max_alt INTEGER NOT NULL,
                                 file_id INTEGER REFERENCES files(id),
                                 UNIQUE(sid,atd,dep_coordinate,arr_coordinate)
                                 );

CREATE TABLE If NOT EXISTS flight_coordinates(
                                                 id SERIAL PRIMARY KEY ,
                                                 sid VARCHAR(100) REFERENCES messages(sid) ON DELETE CASCADE ,
                                                 coordinate geography(POINT, 4326)
);


ALTER TABLE  district_shapes  ADD COLUMN  area_km2 NUMERIC;

UPDATE district_shapes SET area_km2 =  23832      WHERE gid = 1;  -- Сумская обл.
UPDATE district_shapes SET area_km2 = 464275  WHERE gid = 2;  -- Камчатский край
UPDATE district_shapes SET area_km2 = 144902  WHERE gid = 3;  -- Мурманская обл.
UPDATE district_shapes SET area_km2 = 160236  WHERE gid = 4;  -- Пермский край
UPDATE district_shapes SET area_km2 = 194307  WHERE gid = 5;  -- Свердловская обл.
UPDATE district_shapes SET area_km2 = 144527  WHERE gid = 6;  -- Вологодская обл.
UPDATE district_shapes SET area_km2 = 29084   WHERE gid = 7;  -- Владимирская обл.
UPDATE district_shapes SET area_km2 = 15125   WHERE gid = 8;  -- Калининградская обл.
UPDATE district_shapes SET area_km2 = 88529   WHERE gid = 9;  -- Челябинская обл.
UPDATE district_shapes SET area_km2 = 54501   WHERE gid = 10; -- Новгородская обл.
UPDATE district_shapes SET area_km2 = 26081   WHERE gid = 11; -- Республика Крым
UPDATE district_shapes SET area_km2 = 864     WHERE gid = 12; -- Севастополь
UPDATE district_shapes SET area_km2 = 864     WHERE gid = 13; -- Севастополь (дубликат)
UPDATE district_shapes SET area_km2 = NULL    WHERE gid = 14; -- Автономная Республика Крым
UPDATE district_shapes SET area_km2 = 3123    WHERE gid = 15; -- Ингушетия
UPDATE district_shapes SET area_km2 = 84201   WHERE gid = 16; -- Тверская обл.
UPDATE district_shapes SET area_km2 = 87101   WHERE gid = 17; -- Сахалинская обл.
UPDATE district_shapes SET area_km2 = 29777   WHERE gid = 18; -- Калужская обл.
UPDATE district_shapes SET area_km2 = 53565   WHERE gid = 19; -- Самарская обл.
UPDATE district_shapes SET area_km2 = 21437   WHERE gid = 20; -- Ивановская обл.
UPDATE district_shapes SET area_km2 = 24652   WHERE gid = 21; -- Орловская обл.
UPDATE district_shapes SET area_km2 = 49779   WHERE gid = 22; -- Смоленская обл.
UPDATE district_shapes SET area_km2 = 25679   WHERE gid = 23; -- Тульская обл.
UPDATE district_shapes SET area_km2 = 721481  WHERE gid = 24; -- Чукотский АО
UPDATE district_shapes SET area_km2 = 314391  WHERE gid = 25; -- Томская обл.
UPDATE district_shapes SET area_km2 = 164673  WHERE gid = 26; -- Приморский край
UPDATE district_shapes SET area_km2 = 180520  WHERE gid = 27; -- Республика Карелия
UPDATE district_shapes SET area_km2 = 589913  WHERE gid = 28; -- Архангельская обл. (вкл. НАО)
UPDATE district_shapes SET area_km2 = 26128   WHERE gid = 29; -- Мордовия
UPDATE district_shapes SET area_km2 = 37181   WHERE gid = 30; -- Ульяновская обл.
UPDATE district_shapes SET area_km2 = 112877  WHERE gid = 31; -- Волгоградская обл.
UPDATE district_shapes SET area_km2 = 49024   WHERE gid = 32; -- Астраханская обл.
UPDATE district_shapes SET area_km2 = 29997   WHERE gid = 33; -- Курская обл.
UPDATE district_shapes SET area_km2 = 52216   WHERE gid = 34; -- Воронежская обл.
UPDATE district_shapes SET area_km2 = 36177   WHERE gid = 35; -- Ярославская обл.
UPDATE district_shapes SET area_km2 = 177756  WHERE gid = 36; -- Новосибирская обл.
UPDATE district_shapes SET area_km2 = 176810  WHERE gid = 37; -- Ненецкий АО
UPDATE district_shapes SET area_km2 = 416774  WHERE gid = 38; -- Республика Коми
UPDATE district_shapes SET area_km2 = 141140  WHERE gid = 39; -- Омская обл.
UPDATE district_shapes SET area_km2 = 142947  WHERE gid = 40; -- Башкортостан
UPDATE district_shapes SET area_km2 = 123702  WHERE gid = 41; -- Оренбургская обл.
UPDATE district_shapes SET area_km2 = 36271   WHERE gid = 42; -- Еврейская АО
UPDATE district_shapes SET area_km2 = 42061   WHERE gid = 43; -- Удмуртия
UPDATE district_shapes SET area_km2 = 67847   WHERE gid = 44; -- Татарстан
UPDATE district_shapes SET area_km2 = 74731   WHERE gid = 45; -- Калмыкия
UPDATE district_shapes SET area_km2 = 1403    WHERE gid = 46; -- Санкт-Петербург
UPDATE district_shapes SET area_km2 = 76624   WHERE gid = 47; -- Нижегородская обл.
UPDATE district_shapes SET area_km2 = 83908   WHERE gid = 48; -- Ленинградская обл.
UPDATE district_shapes SET area_km2 = 120374  WHERE gid = 49; -- Кировская обл.
UPDATE district_shapes SET area_km2 = 60211   WHERE gid = 50; -- Костромская обл.
UPDATE district_shapes SET area_km2 = 34857   WHERE gid = 51; -- Брянская обл.
UPDATE district_shapes SET area_km2 = 55399   WHERE gid = 52; -- Псковская обл.
UPDATE district_shapes SET area_km2 = 101240  WHERE gid = 53; -- Саратовская обл.
UPDATE district_shapes SET area_km2 = 43352   WHERE gid = 54; -- Пензенская обл.
UPDATE district_shapes SET area_km2 = 24047   WHERE gid = 55; -- Липецкая обл.
UPDATE district_shapes SET area_km2 = 361908  WHERE gid = 56; -- Амурская обл.
UPDATE district_shapes SET area_km2 = 7987    WHERE gid = 57; -- Северная Осетия - Алания
UPDATE district_shapes SET area_km2 = 50270   WHERE gid = 58; -- Дагестан
UPDATE district_shapes SET area_km2 = 16171   WHERE gid = 59; -- Чечня
UPDATE district_shapes SET area_km2 = 787633  WHERE gid = 60; -- Хабаровский край
UPDATE district_shapes SET area_km2 = 462464  WHERE gid = 61; -- Магаданская обл.
UPDATE district_shapes SET area_km2 = 769250  WHERE gid = 62; -- Ямало-Ненецкий АО
UPDATE district_shapes SET area_km2 = 534801  WHERE gid = 63; -- Ханты-Мансийский АО - Югра
UPDATE district_shapes SET area_km2 = 44329   WHERE gid = 64; -- Московская обл.
UPDATE district_shapes SET area_km2 = 2561    WHERE gid = 65; -- Москва
UPDATE district_shapes SET area_km2 = 66160   WHERE gid = 66; -- Ставропольский край
UPDATE district_shapes SET area_km2 = 100967  WHERE gid = 67; -- Ростовская обл.
UPDATE district_shapes SET area_km2 = 75485   WHERE gid = 68; -- Краснодарский край
UPDATE district_shapes SET area_km2 = 7792    WHERE gid = 69; -- Адыгея
UPDATE district_shapes SET area_km2 = 34462   WHERE gid = 70; -- Тамбовская обл.
UPDATE district_shapes SET area_km2 = 39605   WHERE gid = 71; -- Рязанская обл.
UPDATE district_shapes SET area_km2 = 95725   WHERE gid = 72; -- Кемеровская обл.
UPDATE district_shapes SET area_km2 = 61569   WHERE gid = 73; -- Республика Хакасия
UPDATE district_shapes SET area_km2 = 167996  WHERE gid = 74; -- Алтайский край
UPDATE district_shapes SET area_km2 = 92903   WHERE gid = 75; -- Республика Алтай
UPDATE district_shapes SET area_km2 = 431892  WHERE gid = 76; -- Забайкальский край
UPDATE district_shapes SET area_km2 = 23375   WHERE gid = 77; -- Марий Эл
UPDATE district_shapes SET area_km2 = 18343   WHERE gid = 78; -- Чувашия
UPDATE district_shapes SET area_km2 = 27134   WHERE gid = 79; -- Белгородская обл.
UPDATE district_shapes SET area_km2 = 3083523 WHERE gid = 80; -- Республика Саха (Якутия)
UPDATE district_shapes SET area_km2 = 71488   WHERE gid = 81; -- Курганская обл.
UPDATE district_shapes SET area_km2 = 160122  WHERE gid = 82; -- Тюменская обл. (без ХМАО и ЯНАО)
UPDATE district_shapes SET area_km2 = 168604  WHERE gid = 83; -- Тыва
UPDATE district_shapes SET area_km2 = 2366797 WHERE gid = 84; -- Красноярский край
UPDATE district_shapes SET area_km2 = 351334  WHERE gid = 85; -- Бурятия
UPDATE district_shapes SET area_km2 = 774846  WHERE gid = 86; -- Иркутская обл.
UPDATE district_shapes SET area_km2 = 12470   WHERE gid = 87; -- Кабардино-Балкария
UPDATE district_shapes SET area_km2 = 14277   WHERE gid = 88; -- Карачаево-Черкесия
;
