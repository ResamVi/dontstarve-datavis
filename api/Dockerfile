FROM rust:1.31

WORKDIR /usr/src/myapp
COPY Cargo.toml .
COPY Rocket.toml .
COPY src src

RUN rustup default nightly
RUN cargo build --release

CMD ["./target/release/datavizapi"]