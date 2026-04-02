# Hello World

<StoryHeader
    title="First Operation"
    duration="2"
    difficulty="beginner"
/>

## Objective

Run your first phenotype-go-kit operation.

## Implementation

```rust
use phenotype_go_kit::Client;

#[tokio::main]
async fn main() {
    let client = Client::new().await.unwrap();
    let result = client.hello().await.unwrap();
    println!("{}", result);
}
```

## Output

```
Hello from phenotype-go-kit!
```
