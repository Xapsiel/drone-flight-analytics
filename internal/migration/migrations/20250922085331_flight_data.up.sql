
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
                                                 sid VARCHAR(100) REFERENCES messages(sid),
                                                 coordinate geography(POINT, 4326)
);
